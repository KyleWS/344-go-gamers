package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis"
	"github.com/info344-a17/group-alexsirr/servers/gateway/game"
	"github.com/info344-a17/group-alexsirr/servers/gateway/handlers"
	"github.com/info344-a17/group-alexsirr/servers/gateway/models/users"
	"github.com/info344-a17/group-alexsirr/servers/gateway/sessions"
	mgo "gopkg.in/mgo.v2"
)

func main() {
	addr := os.Getenv("ADDR")
	if len(addr) == 0 {
		addr = ":443"
	}

	tlskey := os.Getenv("TLSKEY")
	tlscert := os.Getenv("TLSCERT")

	if len(tlskey) == 0 || len(tlscert) == 0 {
		log.Fatal("please set TLSKEY and TLSCERT")
	}

	sessionkey := os.Getenv("SESSIONKEY")
	if len(sessionkey) == 0 {
		log.Fatal("please set SESSIONKEY")
	}

	raddr := os.Getenv("REDISADDR")
	if len(raddr) == 0 {
		log.Fatal("please set REDISADDR")
	}

	dbaddr := os.Getenv("DBADDR")
	if len(dbaddr) == 0 {
		log.Fatal("please set DBADDR")
	}

	rclient := redis.NewClient(&redis.Options{
		Addr: raddr,
	})
	rstore := sessions.NewRedisStore(rclient, time.Hour)

	mongoSess, err := mgo.Dial(dbaddr)
	if err != nil {
		log.Fatal("Could not dial mongodb")
	}

	dbstore := users.NewMongoStore(mongoSess, "users", "users")

	gs := game.CreateGameState()

	gcontext := game.GameContext{
		UsersStore: dbstore,
		GameState:  gs,
	}

	go gcontext.GameController()

	context := handlers.HandlerContext{
		SigningKey:   sessionkey,
		SessionStore: rstore,
		UsersStore:   dbstore,
		GameState:    gs,
	}

	mux := http.NewServeMux()

	// user and session related handlers
	mux.HandleFunc("/v1/users", context.UsersHandler)
	mux.HandleFunc("/v1/users/me", context.UsersMeHandler)
	mux.HandleFunc("/v1/sessions", context.SessionsHandler)
	mux.HandleFunc("/v1/sessions/mine", context.SessionsMineHandler)

	// game related handlers
	mux.Handle("/v1/game/ws", context.NewGameWebSocketsHandler())
	mux.HandleFunc("/v1/game/submit", context.UploadHandler)

	// File Server
	fs := http.FileServer(http.Dir("./tmp"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	wrappedmux := handlers.NewCORSHandler(mux)

	if os.Getenv("ENV") == "dev" {
		fmt.Printf("server is listening at http://%s\n", addr)
		log.Fatal(http.ListenAndServe(addr, wrappedmux))
	} else {
		fmt.Printf("server is listening at https://%s\n", addr)
		log.Fatal(http.ListenAndServeTLS(addr, tlscert, tlskey, wrappedmux))
	}
}
