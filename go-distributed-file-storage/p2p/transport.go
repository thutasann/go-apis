package p2p

import "net"

// Peer is an interface the represents the remote node
type Peer interface {
	net.Conn
	Send([]byte) error // Write data to the connection
	CloseStream()      // Close the stream
}

// Transport is anything that can handle the communication
// between the nodes in the network.
// This can be of the form (TCP, UDP, websockets)
type Transport interface {
	Addr() string
	Dial(string) error      // Dial the Address and Connect
	ListenAndAccept() error // Listen And Accept
	Consume() <-chan RPC    // Consume
	Close() error           // Close the Transport
}
