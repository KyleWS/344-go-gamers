package handlers

import (
	"github.com/info344-a17/group-alexsirr/servers/gateway/game"
	"github.com/info344-a17/group-alexsirr/servers/gateway/models/users"
	"github.com/info344-a17/group-alexsirr/servers/gateway/sessions"
)

// HandlerContext will act as a receiver for
// HTTP handlers to allow them to access certain parameters
type HandlerContext struct {
	SigningKey   string
	SessionStore sessions.Store
	UsersStore   users.Store
	GameState    *game.GameState
}
