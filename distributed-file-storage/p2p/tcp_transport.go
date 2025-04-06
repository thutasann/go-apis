package p2p

import (
	"fmt"
	"net"
	"sync"
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

// TCP Transport struct
type TCPTransport struct {
	listenAddress string            // Listen Address
	listener      net.Listener      // Network Listener
	mu            sync.RWMutex      // Mutex that will protect Peer
	peers         map[net.Addr]Peer // Peers
}

// Get New TCP Peer
func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

// Get New TCP Transport
func NewTCPTransport(listenAddr string) *TCPTransport {
	return &TCPTransport{
		listenAddress: listenAddr,
	}
}

// Transport Listen And Accept from (TCP, UDP, websockets, etc.)
func (t *TCPTransport) ListenAndAccept() error {
	var err error
	t.listener, err = net.Listen("tcp", t.listenAddress)
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

		go t.handleConn(conn)
	}
}

// Handle TCP Connection
func (t *TCPTransport) handleConn(conn net.Conn) {
	peer := NewTCPPeer(conn, true)
	fmt.Printf("new incoming connection: %v\n", peer)
}
