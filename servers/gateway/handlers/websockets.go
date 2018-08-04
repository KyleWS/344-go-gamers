package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/info344-a17/group-alexsirr/servers/gateway/sessions"
)

//WebSocketsHandler is a handler for WebSocket upgrade requests
type WebSocketsHandler struct {
	upgrader *websocket.Upgrader
	ctx      *HandlerContext
}

//NewGameWebSocketsHandler constructs a new WebSocketsHandler for the game
func (ctx *HandlerContext) NewGameWebSocketsHandler() *WebSocketsHandler {
	return &WebSocketsHandler{
		upgrader: &websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		ctx: ctx,
	}
}

//ServeHTTP implements the http.Handler interface for the WebSocketsHandler
func (wsh *WebSocketsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s := &SessionState{}
	if _, err := sessions.GetState(r, wsh.ctx.SigningKey, wsh.ctx.SessionStore, s); err != nil {
		http.Error(w, "Error getting session: "+err.Error(), http.StatusUnauthorized)
		return
	}

	conn, err := wsh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error upgrading connection %v", err)
		return
	}

	if err := wsh.ctx.GameState.AddClient(conn, s.User.ID); err != nil {
		cm := websocket.FormatCloseMessage(websocket.CloseInternalServerErr, fmt.Sprintf("error adding client: %v", err))
		// no need to check for err, we will be closing the conn regardless
		if err := conn.WriteMessage(websocket.CloseMessage, cm); err != nil {
			fmt.Println(err.Error())
		}
		conn.Close()
	}
}
