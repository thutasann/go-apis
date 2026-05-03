package network

// NetAddr represents the network address of a peer node.
type NetAddr string

// RPC represents a remote procedure call message exchanged between peers.
type RPC struct {
	// From is the network address of the sender peer.
	From NetAddr
	// Payload contains the serialized data of the RPC message.
	Payload []byte
}

// Transport defines the interface for network communication between peers.
type Transport interface {
	// Consume returns a read-only channel that receives incoming RPC messages.
	Consume() <-chan RPC
	// Connect establishes a connection to another peer transport.
	Connect(Transport) error
	// SendMessage sends a message payload to a peer at the specified address.
	SendMessage(NetAddr, []byte) error
	// Addr returns the network address of this transport peer.
	Addr() NetAddr
}
