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

// redis server
type Server struct {
	Config              // Config
	ln     net.Listener // Net Listener
}

// initialize new server
func NewServer(cfg Config) *Server {
	if len(cfg.ListenAddr) == 0 {
		cfg.ListenAddr = defaultListenAddr
	}
	return &Server{
		Config: cfg,
	}
}

// start the server
func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return err
	}
	s.ln = ln
	return s.acceptLoop()
}

// Accept Loop
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

}

func main() {
	fmt.Println("REDIS FROM SCRATCH")
}
