package station

import (
	"bufio"
	"context"
	"log"
	"net"

	"github.com/thutasann/ctp/internal/telemetry"
)

// Station Server accepts train connections
type Server struct {
	Addr   string
	Buffer *Buffer[telemetry.Telemetry]
	Logger *log.Logger
}

func (s *Server) Listen(ctx context.Context) error {
	ln, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}
	defer ln.Close()

	s.Logger.Println("listening on", s.Addr)

	go func() {
		<-ctx.Done()
		ln.Close()
	}()

	for {
		conn, err := ln.Accept()
		if err != nil {
			return nil
		}

		go s.handleConn(ctx, conn)
	}

}

func (s *Server) handleConn(ctx context.Context, conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		t, err := telemetry.Decode(reader)
		if err != nil {
			return
		}

		ok := s.Buffer.Enqueue(ctx, t)
		if !ok {
			return
		}
	}
}
