package control

import (
	"bufio"
	"context"
	"net"

	"github.com/thutasann/ctp/internal/telemetry"
	"github.com/thutasann/ctp/pkg/logx"
)

// Server handles incoming station connections
type Server struct {
	Addr       string
	Aggregator *Aggregator
}

// Listen the incoming connections
func (s *Server) Listen(ctx context.Context) error {
	ln, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}
	defer ln.Close()

	logger := logx.New("CONTROL")
	logger.Println("listening on", s.Addr)

	go func() {
		<-ctx.Done()
		logger.Println("shutdown signal received")
		ln.Close()
	}()

	for {
		conn, err := ln.Accept()
		if err != nil {
			return nil
		}

		go s.handleConn(conn)
	}

}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		t, err := telemetry.Decode(reader)
		if err != nil {
			return
		}

		s.Aggregator.Add(t)
	}
}
