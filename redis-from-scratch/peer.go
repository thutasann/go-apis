package main

import "net"

// Peer struct
type Peer struct {
	conn net.Conn // Peer connection
}

// Initialize new peer
func NewPeer(conn net.Conn) *Peer {
	return &Peer{
		conn: conn,
	}
}

// Read Peer Loop
func (p *Peer) readLoop() {

}
