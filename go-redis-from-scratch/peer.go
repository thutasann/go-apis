package main

import (
	"log/slog"
	"net"
)

// Peer struct
type Peer struct {
	conn  net.Conn     // Peer connection
	msgCh chan Message // Message Channel
}

// Initialize new peer
func NewPeer(conn net.Conn, msgCh chan Message) *Peer {
	return &Peer{
		conn:  conn,
		msgCh: msgCh,
	}
}

// Peer Send Message
func (p *Peer) Send(msg []byte) (int, error) {
	return p.conn.Write(msg)
}

// Read Peer Loop
func (p *Peer) readLoop() error {
	buf := make([]byte, 1024)

	for {
		n, err := p.conn.Read(buf)
		if err != nil {
			slog.Error("peer read error", "error", err)
			return err
		}
		msgBuf := make([]byte, n)
		copy(msgBuf, buf[:n])

		p.msgCh <- Message{
			data: msgBuf,
			peer: p,
		}
	}
}
