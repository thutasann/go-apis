/*
1. A new TCP connection is accepted.

2. A TCPPeer is created for this connection.

3. The HandshakeFunc is called to validate/authenticate the peer.

4. If the handshake is successful:

5. The OnPeer callback is triggered.

6. A loop is started to decode and handle incoming messages.

7. Messages are passed into the rpcch channel.

# Real Life Example

- You’re building a clubhouse (TCPTransport)

- People walk in the door (Accept())

- A bouncer (HandshakeFunc) checks their ID

- If approved, they’re given a name tag (TCPPeer)

- Then they’re handed to a host (OnPeer) to join the party (your network)
*/
package p2p

import (
	"errors"
	"fmt"
	"log"
	"net"
)

// TCPPeer represents the remote node over a TCP estalished connection
type TCPPeer struct {
	// the underlying connection of the peer. which is this case is a TCP connection
	net.Conn

	// meta data to consider :
	//
	// if we dial and retrieve a connection => outbound == true
	//
	// if we accept and retrieve a connection ==> outbound == false
	outbound bool
}

// TCP Transport Options
type TCPTransportOpts struct {
	ListenAddr    string           // Listen Address
	HandshakeFunc HandshakeFunc    // A custom handshake function to validate a peer when a new connection is made.
	Decoder       Decoder          // Decoder
	OnPeer        func(Peer) error // A callback that gets executed when a new peer is successfully connected and handshaked.
}

// TCP Transport struct
type TCPTransport struct {
	TCPTransportOpts              // TCP Transport Options
	listener         net.Listener // TCP listener created with net.Listen.
	rpcch            chan RPC     // Channel used to pass decoded RPC messages to the rest of the system.
}

// Get New TCP Peer
func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		Conn:     conn,
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

// Send implements Peer interface to write data to the connections
func (p *TCPPeer) Send(b []byte) error {
	_, err := p.Conn.Write(b)
	return err
}

// Consume implements the Transport interface which will return read-only channel
// for reading the incoming messages received from another peer in the network
func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcch
}

// ListenAndAccept implements Transport interface
// Transport Listen And Accept from (TCP, UDP, websockets, etc.)
func (t *TCPTransport) ListenAndAccept() error {
	var err error
	t.listener, err = net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}

	go t.startAcceptLoop()

	log.Printf("[ListenAndAccept] TCP transport listening on port: %s\n", t.ListenAddr)

	return nil
}

// Close imlpements the Transport interface
func (t *TCPTransport) Close() error {
	return t.listener.Close()
}

// Dial implements Transport interface
func (t *TCPTransport) Dial(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	go t.handleConn(conn, true)
	return nil
}

// Start TCP Accept Loop
func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()

		if errors.Is(err, net.ErrClosed) {
			return
		}

		if err != nil {
			fmt.Printf("TCP accept error: %s\n", err)
		}

		// fmt.Printf("new incoming connection %+v\n", conn)

		go t.handleConn(conn, false)
	}
}

// Handle TCP Connection
// - Initialize TCP Peer First
// - HandShake
// - Run OnPeer Function
// - Decode The message
func (t *TCPTransport) handleConn(conn net.Conn, outbound bool) {
	var err error

	defer func() {
		fmt.Printf("dropping peer connection: %s", err)
	}()

	// Initialize TCP Peer
	peer := NewTCPPeer(conn, outbound)

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
		err = t.Decoder.Decode(conn, &rpc)

		// todo: handle error properly
		if err != nil {
			fmt.Printf("TCP read Error: %s\n", err)
			return
		}

		rpc.From = conn.RemoteAddr()
		t.rpcch <- rpc
	}

}
