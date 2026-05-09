package main

import (
	"fmt"
	"log/slog"
	"net"
)

// Default Listen Address
const defaultListenAddr = ":5001"

// config struct
type Config struct {
	ListenAddr string // Net Listen Address
}

// Mesasge Struct
type Message struct {
	data []byte // Mesasge Data
	peer *Peer  // Message Peer
}

// redis server
type Server struct {
	Config                   // Config
	peers     map[*Peer]bool // Peers map represents a connected client or node
	ln        net.Listener   // Net Listener
	addPeerCh chan *Peer     // Add Peer Channel
	quitCh    chan struct{}  // Quit channel
	msgCh     chan Message   // Message Channel
	kv        *KV            // Key Value struct
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
		msgCh:     make(chan Message),
		kv:        NewKV(),
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
		case msg := <-s.msgCh: // rawMsg <- from peer.go
			if err := s.handleMesasge(msg); err != nil {
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
	if err := peer.readLoop(); err != nil {
		slog.Error("peer read error", "error", err, "remoteAddr", conn.RemoteAddr())
	}
}

// handle incoming raw message
func (s *Server) handleMesasge(msg Message) error {
	cmd, err := parseCommand(string(msg.data))
	if err != nil {
		return err
	}

	switch v := cmd.(type) {
	case SetCommand:
		return s.kv.Set(v.key, v.val)
	case GetCommand:
		val, ok := s.kv.Get(v.key)
		if !ok {
			return fmt.Errorf("key not found")
		}
		_, err := msg.peer.Send(val)
		if err != nil {
			slog.Error("peer send error", "error", err)
		}
	}

	return nil
}
