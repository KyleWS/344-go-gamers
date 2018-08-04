package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/info344-a17/group-alexsirr/servers/gateway/sessions"

	"github.com/info344-a17/group-alexsirr/servers/gateway/models/users"
)

// UsersHandler will handle new registration requests
// Accepts a POST request with new user data. If valid, a new user and session will be created
func (ctx *HandlerContext) UsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var nu users.NewUser
		if err := json.NewDecoder(r.Body).Decode(&nu); err != nil {
			http.Error(w, "Error decoding JSON", http.StatusBadRequest)
			return
		}
		if err := nu.Validate(); err != nil {
			http.Error(w, "Error invalid user: "+err.Error(), http.StatusBadRequest)
			return
		}
		// if a nil err is returned, that means a user was found
		if _, err := ctx.UsersStore.GetByEmail(nu.Email); err == nil {
			http.Error(w, "Error user already exists with email: "+nu.Email, http.StatusBadRequest)
			return
		}
		if _, err := ctx.UsersStore.GetByUserName(nu.UserName); err == nil {
			http.Error(w, "Error user already exists with username: "+nu.UserName, http.StatusBadRequest)
			return
		}
		u, err := ctx.UsersStore.Insert(&nu)
		if err != nil {
			http.Error(w, "Error creating user", http.StatusInternalServerError)
			return
		}
		ss := SessionState{
			SessionStart: time.Now(),
			User:         *u,
		}
		if _, err := sessions.BeginSession(ctx.SigningKey, ctx.SessionStore, ss, w); err != nil {
			http.Error(w, "Error creating session", http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(u); err != nil {
			http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
			return
		}

	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

// helper function to get the intersection of two slices of bson ids
func getBidsIntersect(bids1 []bson.ObjectId, bids2 []bson.ObjectId) []bson.ObjectId {
	res := []bson.ObjectId{}
	seen := make(map[bson.ObjectId]bool)
	// set seen for all items in bids1
	for _, bid := range bids1 {
		seen[bid] = true
	}
	// if an item was seen in bids1 and bids2, add to res
	for _, bid := range bids2 {
		if seen[bid] {
			res = append(res, bid)
		}
	}

	return res
}

// UsersMeHandler will handle a user that is currently signed in
// GET requests will return the session based on the Authorization header
// PATCH will update a user
func (ctx *HandlerContext) UsersMeHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		s := &SessionState{}
		if _, err := sessions.GetState(r, ctx.SigningKey, ctx.SessionStore, s); err != nil {
			http.Error(w, "Error getting session: "+err.Error(), http.StatusUnauthorized)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(s.User); err != nil {
			http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
			return
		}

	case "PATCH":
		s := &SessionState{}
		id, err := sessions.GetState(r, ctx.SigningKey, ctx.SessionStore, s)
		if err != nil {
			http.Error(w, "Error getting session: "+err.Error(), http.StatusUnauthorized)
			return
		}

		up := users.Updates{}
		if err := json.NewDecoder(r.Body).Decode(&up); err != nil {
			http.Error(w, "Error decoding JSON", http.StatusBadRequest)
			return
		}
		// update user in UsersStore
		if err := ctx.UsersStore.Update(s.User.ID, &up); err != nil {
			http.Error(w, "Error updating user: "+err.Error(), http.StatusBadRequest)
			return
		}
		// apply updates to user to update in state
		if err := s.User.ApplyUpdates(&up); err != nil {
			http.Error(w, "Error updating user: "+err.Error(), http.StatusBadRequest)
			return
		}
		// update user in SessionStore
		if err := ctx.SessionStore.Save(id, s); err != nil {
			http.Error(w, "Error saving session", http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(s.User); err != nil {
			http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

// SessionsHandler will handle user login via a POST request
func (ctx *HandlerContext) SessionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		creds := users.Credentials{}
		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			http.Error(w, "Error decoding JSON", http.StatusBadRequest)
			return
		}
		u, err := ctx.UsersStore.GetByEmail(creds.Email)
		if err != nil {
			// if email is not found, create a temp user and set password
			// this will prevent timing attacks from taking place
			temp := &users.User{}
			temp.SetPassword("nothing")
			http.Error(w, "Error invalid email/password", http.StatusUnauthorized)
			return
		}
		if err := u.Authenticate(creds.Password); err != nil {
			http.Error(w, "Error invalid email/password", http.StatusUnauthorized)
			return
		}
		ss := SessionState{
			SessionStart: time.Now(),
			User:         *u,
		}
		if _, err := sessions.BeginSession(ctx.SigningKey, ctx.SessionStore, ss, w); err != nil {
			http.Error(w, "Error creating session", http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(u); err != nil {
			http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

// SessionsMineHandler will sign a user out when a DELETE request with valid
// auth is sent
func (ctx *HandlerContext) SessionsMineHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "DELETE" {
		if _, err := sessions.EndSession(r, ctx.SigningKey, ctx.SessionStore); err != nil {
			http.Error(w, "Error getting session: "+err.Error(), http.StatusUnauthorized)
			return
		}
		w.Write([]byte("signed out"))
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}
