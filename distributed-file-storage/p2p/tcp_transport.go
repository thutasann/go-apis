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

// TCP Transport Options
type TCPTransportOpts struct {
	ListenAddr    string           // Listen Address
	HandshakeFunc HandshakeFunc    // HandShake Function
	Decoder       Decoder          // Decoder
	OnPeer        func(Peer) error // On Peer Function
}

// TCP Transport struct
type TCPTransport struct {
	TCPTransportOpts              // TCP Transport Options
	listener         net.Listener // Net Listener
	rpcch            chan RPC     // RPC Channel
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
		rpcch:            make(chan RPC),
	}
}

// Close implements the Peer interface.
func (p *TCPPeer) Close() error {
	return p.conn.Close()
}

// Consume implements the `Transport` interface which will return read-only channel
// for reading the incoming messages received from another peer in the network
func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcch
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
// - Initialize TCP Peer First
// - HandShake
// - Decode The message
func (t *TCPTransport) handleConn(conn net.Conn) {
	var err error

	defer func() {
		fmt.Printf("dropping peer connection: %s", err)
	}()

	// Initialize TCP Peer
	peer := NewTCPPeer(conn, true)

	// HandShake
	if err := t.HandshakeFunc(peer); err != nil {
		conn.Close()
		fmt.Printf("TCP HandShake Error: %s\n", err)
		return
	}

	if t.OnPeer != nil {
		if err = t.OnPeer(peer); err != nil {
			return
		}
	}

	// Read loop
	rpc := RPC{}
	for {
		if err := t.Decoder.Decode(conn, &rpc); err != nil {
			fmt.Printf("TCP Error: %s\n", err)
			continue
		}

		rpc.From = conn.RemoteAddr()
		t.rpcch <- rpc
	}

}
