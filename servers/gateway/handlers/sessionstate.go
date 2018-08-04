package handlers

import (
	"time"

	"github.com/info344-a17/group-alexsirr/servers/gateway/models/users"
)

// SessionState tracks information about a session
type SessionState struct {
	SessionStart time.Time
	User         users.User
}
