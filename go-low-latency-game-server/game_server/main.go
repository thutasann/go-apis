package main

import (
	"fmt"
	"net/http"

	"github.com/anthdm/hollywood/actor"
	"github.com/gorilla/websocket"
)

type GameServer struct{}

func newGameServer() actor.Receiver {
	return &GameServer{}
}

func (s *GameServer) Receive(c *actor.Context) {
	switch msg := c.Message().(type) {
	case actor.Started:
		s.startHTTP()
		_ = msg
	}
}

func (s *GameServer) startHTTP() {
	fmt.Println("starting HTTP server on port 40000")
	go func() {
		http.HandleFunc("/ws", s.handleWS)
		http.ListenAndServe(":40000", nil)
	}()
}

// handles the upgrade of the websocekt
func (s *GameServer) handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if err != nil {
		fmt.Println("ws upgrade error: ", err)
		return
	}
	fmt.Println("new client trying connect")
	fmt.Println(conn)
}

// game_server
func main() {
	e, _ := actor.NewEngine(actor.NewEngineConfig())
	e.Spawn(newGameServer, "server")
	select {}
}
