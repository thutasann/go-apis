package p2p

// Peer is an interface the represents the remote node
type Peer interface {
	Close() error // Peer Close Function
}

// Transport is anything that can handle the communication
// between the nodes in the network.
// This can be of the form (TCP, UDP, websockets)
type Transport interface {
	ListenAndAccept() error // Listen And Accept
	Consume() <-chan RPC    // Consume
}
