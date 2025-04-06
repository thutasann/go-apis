package p2p

import (
	"fmt"
	"net"
)

// TCPPeer represents the remote node over a TCP estalished connection
type TCPPeer struct {
	// conn is the underlying connection of the peer
	conn net.Conn

	// if we dial and retrieve a connection => outbound == true
	//
	// if we accept and retrieve a connection ==> outbound == false
	outbound bool
}

type TCPTransportOpts struct {
	ListenAddr    string        // Listen Address
	HandshakeFunc HandshakeFunc // HandShake Function
	Decoder       Decoder       // Decoder
}

// TCP Transport struct
type TCPTransport struct {
	TCPTransportOpts              // TCP Transport Options
	listener         net.Listener // Net Listener
	// mu               sync.RWMutex      // Mutex that will protect Peer
	// peers            map[net.Addr]Peer // Peers
}

// Get New TCP Peer
func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

// Get New TCP Transport
func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
	}
}

// Transport Listen And Accept from (TCP, UDP, websockets, etc.)
func (t *TCPTransport) ListenAndAccept() error {
	var err error
	t.listener, err = net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}

	go t.startAcceptLoop()
	return nil
}

// Start TCP Accept Loop
func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("TCP accept error: %s\n", err)
		}

		fmt.Printf("new incoming connection %+v\n", conn)

		go t.handleConn(conn)
	}
}

// Handle TCP Connection
// HandShake First
// Decode The message
func (t *TCPTransport) handleConn(conn net.Conn) {
	peer := NewTCPPeer(conn, true)

	if err := t.HandshakeFunc(peer); err != nil {
		conn.Close()
		fmt.Printf("TCP HandShake Error: %s\n", err)
		return
	}

	// Read loop
	msg := &struct{}{}
	for {
		if err := t.Decoder.Decode(conn, msg); err != nil {
			fmt.Printf("TCP Error: %s\n", err)
			continue
		}
	}

}
