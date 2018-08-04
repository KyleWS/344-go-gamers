package sessions

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

const headerAuthorization = "Authorization"
const paramAuthorization = "auth"
const schemeBearer = "Bearer "

//ErrNoSessionID is used when no session ID was found in the Authorization header
var ErrNoSessionID = errors.New("no session ID found in " + headerAuthorization + " header")

//ErrInvalidScheme is used when the authorization scheme is not supported
var ErrInvalidScheme = errors.New("authorization scheme not supported")

//BeginSession creates a new SessionID, saves the `sessionState` to the store, adds an
//Authorization header to the response with the SessionID, and returns the new SessionID
func BeginSession(signingKey string, store Store, sessionState interface{}, w http.ResponseWriter) (SessionID, error) {
	id, err := NewSessionID(signingKey)
	if err != nil {
		return InvalidSessionID, fmt.Errorf("error creating session key: %v", err)
	}
	if err := store.Save(id, sessionState); err != nil {
		return InvalidSessionID, fmt.Errorf("error saving session: %v", err)
	}
	w.Header().Add(headerAuthorization, schemeBearer+id.String())

	return id, nil
}

//GetSessionID extracts and validates the SessionID from the request headers
func GetSessionID(r *http.Request, signingKey string) (SessionID, error) {
	auth := r.Header.Get(headerAuthorization)

	if len(auth) == 0 {
		auth = r.URL.Query().Get("auth")
		if len(auth) == 0 {
			return InvalidSessionID, ErrNoSessionID
		}
	}

	if !strings.HasPrefix(auth, schemeBearer) {
		return InvalidSessionID, ErrInvalidScheme
	}

	id, err := ValidateID(auth[len(schemeBearer):], signingKey)
	if err != nil {
		return InvalidSessionID, err
	}
	return id, nil
}

//GetState extracts the SessionID from the request,
//gets the associated state from the provided store into
//the `sessionState` parameter, and returns the SessionID
func GetState(r *http.Request, signingKey string, store Store, sessionState interface{}) (SessionID, error) {
	id, err := GetSessionID(r, signingKey)
	if err != nil {
		return InvalidSessionID, err
	}
	if err := store.Get(id, sessionState); err != nil {
		return InvalidSessionID, ErrStateNotFound
	}
	return id, nil
}

//EndSession extracts the SessionID from the request,
//and deletes the associated data in the provided store, returning
//the extracted SessionID.
func EndSession(r *http.Request, signingKey string, store Store) (SessionID, error) {
	id, err := GetSessionID(r, signingKey)
	if err != nil {
		return InvalidSessionID, err
	}
	store.Delete(id)
	return id, nil
}
