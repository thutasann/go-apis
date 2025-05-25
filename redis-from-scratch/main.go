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
	quitCh    chan struct{}  // Quit channel
	msgCh     chan []byte    // Message Channel
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
		quitCh:    make(chan struct{}),
		msgCh:     make(chan []byte),
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

	slog.Info("Server running", "ListenAddr", s.ListenAddr)

	return s.acceptLoop()
}

// loop the server
func (s *Server) loop() {
	for {
		select {
		case rawMsg := <-s.msgCh: // rawMsg <- from peer.go
			if err := s.handleRawMesasge(rawMsg); err != nil {
				slog.Error("handle raw message error", "error", err)
			}
		case <-s.quitCh:
			fmt.Println("::: quit channel :::")
			return
		case peer := <-s.addPeerCh:
			s.peers[peer] = true
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
	peer := NewPeer(conn, s.msgCh)
	s.addPeerCh <- peer
	slog.Info("new peer connected", "remoteAddr", conn.RemoteAddr())
	if err := peer.readLoop(); err != nil {
		slog.Error("peer read error", "error", err, "remoteAddr", conn.RemoteAddr())
	}
}

// handle incoming raw message
func (s *Server) handleRawMesasge(rawMsg []byte) error {
	fmt.Println("rawMsg :>> ", string(rawMsg))
	return nil
}

// REDIS FROM SCRATCH
func main() {
	fmt.Println("REDIS FROM SCRATCH")

	server := NewServer(Config{})
	log.Fatal(server.Start())
}
