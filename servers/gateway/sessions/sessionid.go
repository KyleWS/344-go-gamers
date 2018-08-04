package sessions

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
)

//InvalidSessionID represents an empty, invalid session ID
const InvalidSessionID SessionID = ""

//idLength is the length of the ID portion
const idLength = 32

//signedLength is the full length of the signed session ID
//(ID portion plus signature)
const signedLength = idLength + sha256.Size

//SessionID represents a valid, digitally-signed session ID.
//This is a base64 URL encoded string created from a byte slice
//where the first `idLength` bytes are crytographically random
//bytes representing the unique session ID, and the remaining bytes
//are an HMAC hash of those ID bytes (i.e., a digital signature).
//The byte slice layout is like so:
//+-----------------------------------------------------+
//|...32 crypto random bytes...|HMAC hash of those bytes|
//+-----------------------------------------------------+
type SessionID string

//ErrInvalidID is returned when an invalid session id is passed to ValidateID()
var ErrInvalidID = errors.New("Invalid Session ID")

//NewSessionID creates and returns a new digitally-signed session ID,
//using `signingKey` as the HMAC signing key. An error is returned only
//if there was an error generating random bytes for the session ID
func NewSessionID(signingKey string) (SessionID, error) {
	if len(signingKey) == 0 {
		return InvalidSessionID, errors.New("provided signing key was of length 0")
	}

	randBytes := make([]byte, idLength)
	_, err := rand.Read(randBytes)
	if err != nil {
		return InvalidSessionID, fmt.Errorf("error generating rand bytes: %v", err)
	}

	sig := generateHMAC(signingKey, randBytes)
	// allocate new byte slice of length signedLength
	sID := make([]byte, signedLength)

	// copy randBytes into sID
	copy(sID, randBytes)
	// copy sig to last 32 bytes
	copy(sID[32:], sig)
	return SessionID(base64.URLEncoding.EncodeToString(sID)), nil
}

// returns an HMAC byte slice
// takes a signingKey and data to sign
func generateHMAC(signingKey string, data []byte) []byte {
	h := hmac.New(sha256.New, []byte(signingKey))
	h.Write(data)
	sig := h.Sum(nil)
	return sig
}

//ValidateID validates the string in the `id` parameter
//using the `signingKey` as the HMAC signing key
//and returns an error if invalid, or a SessionID if valid
func ValidateID(id string, signingKey string) (SessionID, error) {
	decoded, err := base64.URLEncoding.DecodeString(id)
	if err != nil {
		return InvalidSessionID, ErrInvalidID
	}

	expectedSig := generateHMAC(signingKey, decoded[:32])

	if subtle.ConstantTimeCompare(expectedSig, decoded[32:]) != 1 {
		return InvalidSessionID, ErrInvalidID
	}

	return SessionID(id), nil
}

//String returns a string representation of the sessionID
func (sid SessionID) String() string {
	return string(sid)
}