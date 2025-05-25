package main

import (
	"fmt"
	"log"
	"log/slog"
	"net"
)

// Default Listen Address
const defaultListenAddr = ":5001"

// config struct
type Config struct {
	ListenAddr string // Net Listen Address
}

// redis server
type Server struct {
	Config                   // Config
	peers     map[*Peer]bool // Peers map represents a connected client or node
	ln        net.Listener   // Net Listener
	addPeerCh chan *Peer     // Add Peer Channel
}

// initialize new server
func NewServer(cfg Config) *Server {
	if len(cfg.ListenAddr) == 0 {
		cfg.ListenAddr = defaultListenAddr
	}
	return &Server{
		Config:    cfg,
		peers:     make(map[*Peer]bool),
		addPeerCh: make(chan *Peer),
	}
}

// start the server
func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return err
	}
	s.ln = ln

	go s.loop()

	return s.acceptLoop()
}

// loop the server
func (s *Server) loop() {
	for {
		select {
		case peer := <-s.addPeerCh:
			s.peers[peer] = true
		default:
			fmt.Println("no loop case")
		}
	}
}

// accept Loop
func (s *Server) acceptLoop() error {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			slog.Error("accept error", "err", err)
			continue
		}
		go s.handleConn(conn)
	}
}

// handle listen connection
func (s *Server) handleConn(conn net.Conn) {
	peer := NewPeer(conn)
	s.addPeerCh <- peer
	peer.readLoop()
}

// REDIS FROM SCRATCH
func main() {
	fmt.Println("REDIS FROM SCRATCH")

	server := NewServer(Config{})
	log.Fatal(server.Start())
}
