package game

import "github.com/info344-a17/group-alexsirr/servers/gateway/models/users"

type GameContext struct {
	UsersStore users.Store
	GameState  *GameState
}
