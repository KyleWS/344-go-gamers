package handlers

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/info344-a17/group-alexsirr/servers/gateway/game"
	"github.com/info344-a17/group-alexsirr/servers/gateway/sessions"
)

// fileDir is the directory on the server that stores the static assets (both saved and retrieved
// from here at the moment). This folder must be created before the server tries to save
// files into the directory.

const FileDir = "./tmp"

// UploadHandler handles client round submissions. If a client submits an image (in JSON),
// the image is hashed and written to disk.
func (ctx *HandlerContext) UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		s := &SessionState{}
		if _, err := sessions.GetState(r, ctx.SigningKey, ctx.SessionStore, s); err != nil {
			http.Error(w, "Error getting session: "+err.Error(), http.StatusUnauthorized)
			return
		}

		buff, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("error reading request body: %v", err), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
		req1 := ioutil.NopCloser(bytes.NewBuffer(buff))
		req2 := ioutil.NopCloser(bytes.NewBuffer(buff))
		req3 := ioutil.NopCloser(bytes.NewBuffer(buff))

		h := sha256.New()
		if _, err := io.Copy(h, req1); err != nil {
			http.Error(w, fmt.Sprintf("error reading request buffer while hashing: %v", err), http.StatusInternalServerError)
			return
		}
		hashedFile := h.Sum(nil)
		encodedHashFileName := fmt.Sprintf("%x", hashedFile)

		var RequestStruct interface{}
		body, err := ioutil.ReadAll(req2)
		if err != nil {
			http.Error(w, fmt.Sprintf("error reading request buffer before unmarshalling: %v", err), http.StatusInternalServerError)
			return
		}
		if err := json.Unmarshal(body, &RequestStruct); err != nil {
			http.Error(w, fmt.Sprintf("error unmarshalling request: %v", err), http.StatusBadRequest)
			return
		}
		interfacedStruct := RequestStruct.(map[string]interface{})
		var rs *game.RoundSubmission
		if interfacedStruct["type"] == "drawing" {
			if _, err := os.Stat(FileDir + "/" + encodedHashFileName); os.IsNotExist(err) {
				f, err := os.OpenFile(FileDir+"/"+encodedHashFileName, os.O_CREATE|os.O_WRONLY, 0666)
				if err != nil {
					http.Error(w, fmt.Sprintf("error opening file before writing to disk: %v", err), http.StatusInternalServerError)
					return
				}
				defer f.Close()
				if _, err := io.Copy(f, req3); err != nil {
					http.Error(w, fmt.Sprintf("error writing to disk: %v", err), http.StatusInternalServerError)
					return
				}
			}
			rs = &game.RoundSubmission{
				Owner:     s.User.ID,
				ImageHash: encodedHashFileName,
			}
		} else {
			rs = &game.RoundSubmission{
				Owner:       s.User.ID,
				Description: interfacedStruct["data"].(string),
			}
		}
		ctx.GameState.IncomingRoundSubmissions <- rs
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(rs)
	} else if r.Method != http.MethodOptions {
		http.Error(w, "Only POST is allowed", http.StatusMethodNotAllowed)
		return
	}
}
