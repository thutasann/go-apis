package main

import (
	"log/slog"
	"net"
)

// Peer struct
type Peer struct {
	conn  net.Conn    // Peer connection
	msgCh chan []byte // Message Channel
}

// Initialize new peer
func NewPeer(conn net.Conn, msgCh chan []byte) *Peer {
	return &Peer{
		conn:  conn,
		msgCh: msgCh,
	}
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

		p.msgCh <- msgBuf
	}
}
