package game

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/info344-a17/group-alexsirr/servers/gateway/models/users"
	"gopkg.in/mgo.v2/bson"
)

// seconds per round
const descriptionTime = 15
const drawingTime = 10
const intermissionTime = 15

const drawingphase = "drawing-phase"
const descriptionphase = "description-phase"

// some game settings
const minplayers = 2
const maxRounds = 9

// GameState will control the state of the entire game
// it holds not only a list of all clients connected, but also
// information about the current and past rounds
type GameState struct {
	gamestarted          bool
	roundNumber          int
	currentRoundDuration int
	Phase                string
	playerCount          int
	clients              []*gameClient
	eventQ               chan *gameEvent
	gameSubmissions      []*gameSubmission
	// channel that the submission handler will use to get
	// incoming round submissions to the game
	IncomingRoundSubmissions chan *RoundSubmission
	start                    chan int
	quit                     chan int
	mx                       *sync.Mutex
}

// RoundSubmission represents the submission after a client
// on round end, Description or ImageHash might be nil depending on the
// round type
type RoundSubmission struct {
	Description string
	ImageHash   string
	Owner       bson.ObjectId
}

// gameSubmission holds an entire "stack" of roundsubmissions
// this is what is presented after the game ends
// desc --> drawing --> desc --> etc
type gameSubmission struct {
	roundSubmission *RoundSubmission
	next            *gameSubmission
}

// gameClient represents a client who is
// currently connected. Being connected does not mean
// they are playing
type gameClient struct {
	client    *websocket.Conn
	userID    bson.ObjectId
	isPlaying bool
}

// gameEvent will be sent to clients
// target is the target client to send to
type gameEvent struct {
	target *websocket.Conn
	data   []byte
}

// gameResult is what will be sent back to the clients
// when the game ends, similar to a RoundSubmission, but
// this contains all the user info rather than just their ID
type gameResult struct {
	Owner       *users.User `json:"owner"`
	ImageHash   string      `json:"imageHash"`
	Description string      `json:"description"`
}

// CreateGameState creates an initial GameState and starts
// the notification process
func CreateGameState() *GameState {
	gs := &GameState{
		gamestarted:          false,
		roundNumber:          0,
		currentRoundDuration: descriptionTime,
		Phase:                descriptionphase,
		clients:              []*gameClient{},
		eventQ:               make(chan *gameEvent),
		gameSubmissions:      make([]*gameSubmission, 0, 0),
		start:                make(chan int, 1),
		quit:                 make(chan int, 1),
		mx:                   new(sync.Mutex),
	}
	go gs.startnotifying()
	return gs
}

// GameController is the control center for the game
// it will start, end , and restart the game
func (ctx *GameContext) GameController() {
	for {
		ctx.GameState.waitForPlayers()
		// create a buffered channel based on number of players in game
		ctx.GameState.IncomingRoundSubmissions = make(chan *RoundSubmission, ctx.GameState.playerCount)
		// get list of all users who are playing based on the connected clients
		ulist := []*users.User{}
		for _, c := range ctx.GameState.clients {
			u, err := ctx.UsersStore.GetByID(c.userID)
			if err != nil {
				// if an error occurs looking up the user, remove the client
				ctx.GameState.removeClient(c)
				continue
			}
			ulist = append(ulist, u)
		}
		gameError := ctx.GameState.startGame(ulist)
		ctx.GameState.endGame(ctx.UsersStore, gameError)
		ctx.GameState.resetGame()
	}
}

// function that will block until minplayers are accounted for
func (gs *GameState) waitForPlayers() {
	// add any players who were not playing in the last round to the game
	for _, c := range gs.clients {
		if !c.isPlaying {
			c.isPlaying = true
			gs.mx.Lock()
			gs.playerCount++
			gs.mx.Unlock()
		}
	}
	if gs.playerCount >= minplayers {
		return
	}
	// using a select loop because it blocks the current routine
	// this avoids 100% cpu consumtion from a for {}
	select {
	case <-gs.start:
		return
	}
}

// startGame begins the game
func (gs *GameState) startGame(ulist []*users.User) error {

	gs.gamestarted = true

	msg, err := json.Marshal(map[string]interface{}{
		"action":  "game-start",
		"players": ulist,
	})
	if err != nil {
		fmt.Printf("error marshalling JSON: %v", err)
		select {
		case gs.quit <- 1:
		default:
		}
	}
	gs.notify(&gameEvent{
		target: nil,
		data:   msg,
	})

	// initial description phase
	if err := gs.handlePhaseOne(); err != nil {
		return err
	}

	// while the game is less than the max rounds or 2 * number of playerCount -1
	for gs.roundNumber < 2*gs.playerCount-1 && gs.roundNumber < maxRounds {
		if err := gs.beginphase(); err != nil {
			return err
		}
		if err := gs.beginphase(); err != nil {
			return err
		}
	}
	return nil
}

// EndGame will end the current game
func (gs *GameState) endGame(uStore users.Store, gameError error) {
	msgInterface := map[string]interface{}{
		"action": "game-over",
	}
	if gameError != nil {
		msgInterface["error"] = gameError.Error()
	}

	msg, err := json.Marshal(msgInterface)
	if err != nil {
		fmt.Printf("error marshalling JSON: %v", err)
		select {
		case gs.quit <- 1:
		default:
		}
	}
	gs.notify(&gameEvent{
		target: nil,
		data:   msg,
	})

	ownerMap := make(map[bson.ObjectId]*users.User)
	results := make([][]*gameResult, len(gs.gameSubmissions))

	// for each of the gameSubmission trees
	for i, gameSub := range gs.gameSubmissions {
		curr := gameSub
		// navigate through each tree, pulling out data
		for curr != nil {
			if curr.roundSubmission == nil {
				break
			}
			// cache user pointers to reduce db lookups
			owner, found := ownerMap[curr.roundSubmission.Owner]
			if !found {
				u, err := uStore.GetByID(curr.roundSubmission.Owner)
				if err != nil {
					fmt.Printf("error marshalling JSON: %v", err)
					select {
					case gs.quit <- 1:
					default:
					}
				}
				ownerMap[u.ID] = u
				owner = u
			}
			results[i] = append(results[i], &gameResult{
				Owner:       owner,
				ImageHash:   curr.roundSubmission.ImageHash,
				Description: curr.roundSubmission.Description,
			})
			curr = curr.next
		}
	}

	// broadcast game results to all clients
	msg, err = json.Marshal(map[string]interface{}{
		"action":  "game-results",
		"results": results,
	})
	if err != nil {
		fmt.Printf("error marshalling JSON: %v", err)
		select {
		case gs.quit <- 1:
		default:
		}
	}
	gs.notify(&gameEvent{
		target: nil,
		data:   msg,
	})
}

// resetGame will put game into intermission and reset
// fields in the GameState
func (gs *GameState) resetGame() {
	msg, err := json.Marshal(map[string]interface{}{
		"action":   "intermission",
		"duration": intermissionTime,
	})
	if err != nil {
		fmt.Printf("error marshalling JSON: %v", err)
		select {
		case gs.quit <- 1:
		default:
		}
	}
	gs.notify(&gameEvent{
		target: nil,
		data:   msg,
	})
	// reset important games state values

	gs.gamestarted = false
	gs.Phase = descriptionphase
	gs.currentRoundDuration = descriptionTime
	gs.roundNumber = 0
	gs.gameSubmissions = make([]*gameSubmission, 0, 0)
	gs.quit = make(chan int, 1)
	gs.start = make(chan int)

	<-time.NewTimer(time.Second * time.Duration(intermissionTime)).C
}

func (gs *GameState) beginphase() error {
	// get the rotation for this round
	roundEvents := gs.setRotation()

	// then send to each client individually
	for _, e := range roundEvents {
		gs.notify(e)
	}

	// start new timer
	//<-time.NewTimer(time.Second * time.Duration(gs.currentRoundDuration)).C

	return gs.roundOver()
}

func (gs *GameState) handlePhaseOne() error {
	for i := 0; i < len(gs.clients); i++ {
		gs.gameSubmissions = append(gs.gameSubmissions, &gameSubmission{})
	}

	msg, err := json.Marshal(map[string]interface{}{
		"action":         "first-phase",
		"round-duration": descriptionTime,
	})
	if err != nil {
		fmt.Printf("error marshalling JSON: %v", err)
		select {
		case gs.quit <- 1:
		default:
		}
	}
	gs.notify(&gameEvent{
		target: nil,
		data:   msg,
	})
	//<-time.NewTimer(time.Second * time.Duration(gs.currentRoundDuration)).C

	return gs.roundOver()
}

func (gs *GameState) roundOver() error {
	// wait for player submissions to come in
	results, err := gs.waitForRoundSubmissions()
	if err != nil {
		return err
	}
	// add the round submissions to the game submissions
	for _, rs := range results {
		gameSub := gs.getRoundSubmissionLoc(rs.Owner)

		if gs.roundNumber == 0 {
			gameSub.roundSubmission = rs
		} else {
			gameSub.next = &gameSubmission{roundSubmission: rs}
		}
	}

	// toggle round type and time
	if gs.Phase == descriptionphase {
		gs.currentRoundDuration = drawingTime
		gs.Phase = drawingphase
	} else {
		gs.currentRoundDuration = descriptionTime
		gs.Phase = descriptionphase
	}

	gs.roundNumber++
	return nil
}

// will block until the correct number of submissions come in
func (gs *GameState) waitForRoundSubmissions() ([]*RoundSubmission, error) {
	results := make([]*RoundSubmission, 0, 0)
	for {
		select {
		case rs := <-gs.IncomingRoundSubmissions:
			results = append(results, rs)
			if len(results) == gs.playerCount {
				return results, nil
			}
		case <-gs.quit:
			return nil, errors.New("Game exiting with error")
		}
	}
}

// returns a pointer to the gameSubmission where the client should append to based
// on the current round number
func (gs *GameState) getRoundSubmissionLoc(owner bson.ObjectId) *gameSubmission {
	var curr *gameSubmission
	for i, c := range gs.clients {
		if c.userID == owner {
			curr = gs.gameSubmissions[(i+gs.roundNumber)%len(gs.gameSubmissions)]
		}
	}
	for i := 0; i < gs.roundNumber-1; i++ {
		curr = curr.next
	}
	return curr
}

func (gs *GameState) setRotation() []*gameEvent {
	roundEvents := []*gameEvent{}
	for i, c := range gs.clients {
		if !c.isPlaying {
			continue
		}
		curr := gs.gameSubmissions[(i+gs.roundNumber)%len(gs.gameSubmissions)]
		for j := 0; j < gs.roundNumber-1; j++ {
			curr = curr.next
		}

		var data string
		if gs.Phase == descriptionphase {
			data = curr.roundSubmission.ImageHash
		} else {
			data = curr.roundSubmission.Description
		}

		msg, err := json.Marshal(map[string]interface{}{
			"action":         gs.Phase,
			"round-duration": gs.currentRoundDuration,
			"data":           data,
		})
		if err != nil {
			fmt.Printf("error marshalling JSON: %v", err)
			select {
			case gs.quit <- 1:
			default:
			}
		}
		roundEvents = append(roundEvents, &gameEvent{
			target: c.client,
			data:   msg,
		})
	}

	return roundEvents
}

// GameState notification and client functions

//notify inserts event into the gamestate eventQ
func (gs *GameState) notify(event *gameEvent) {
	gs.eventQ <- event
}

//start starts the notification loop
func (gs *GameState) startnotifying() {
	for e := range gs.eventQ {
		for _, c := range gs.clients {
			// check if target is broadcast or the specific client
			if c.isPlaying && (e.target == nil || e.target == c.client) {
				if err := c.client.WriteMessage(websocket.TextMessage, e.data); err != nil {
					gs.removeClient(c)
				}
			}
		}
	}
}

// AddClient will check and see if a user already is playing the game
// else the connection will be added to the gamestate
func (gs *GameState) AddClient(client *websocket.Conn, userID bson.ObjectId) error {
	for _, c := range gs.clients {
		if c.userID == userID {
			return errors.New("user already has a game session")
		}
	}
	newClient := &gameClient{
		client: client,
		userID: userID,
	}

	// if the game has not started, set the client to be playing
	if !gs.gamestarted {
		newClient.isPlaying = true
		gs.mx.Lock()
		gs.playerCount++
		gs.mx.Unlock()
	}

	gs.mx.Lock()
	gs.clients = append(gs.clients, newClient)
	gs.mx.Unlock()
	if !gs.gamestarted && len(gs.clients) >= minplayers {
		select {
		// add to start channel unless another routine has done so
		case gs.start <- 0:
		default:
		}

	}
	for {
		if _, _, err := client.NextReader(); err != nil {
			gs.removeClient(newClient)
			break
		}
	}
	return nil
}

// function that will find client in gameClient.clients and remove it
func (gs *GameState) removeClient(client *gameClient) {
	for i := range gs.clients {
		if gs.clients[i] == client {
			gs.mx.Lock()
			copy(gs.clients[i:], gs.clients[i+1:])
			gs.clients[len(gs.clients)-1] = nil
			gs.clients = gs.clients[:len(gs.clients)-1]
			gs.mx.Unlock()
			// if the client was playing, subtract 1 from current playerCount
			if client.isPlaying {
				gs.mx.Lock()
				gs.playerCount--
				gs.mx.Unlock()
				// end the game if a player leaves
				if gs.gamestarted {
					//gs.endGame()
					select {
					// add to quit channel unless another routine has done so
					case gs.quit <- 1:
					default:
					}
				}
			}
			break
		}
	}
	client.client.Close()
}
