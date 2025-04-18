package p2p

import "net"

// Peer is an interface the represents the remote node
type Peer interface {
	Send([]byte) error    // Write data to the connection
	RemoteAddr() net.Addr // the remote address of its underlying connection
	Close() error         // Peer Close Function
}

// Transport is anything that can handle the communication
// between the nodes in the network.
// This can be of the form (TCP, UDP, websockets)
type Transport interface {
	Dial(string) error      // Dial the Address and Connect
	ListenAndAccept() error // Listen And Accept
	Consume() <-chan RPC    // Consume
	Close() error           // Close the Transport
}
