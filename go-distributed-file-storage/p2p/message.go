package p2p

const (
	IncomingMessage = 0x1
	IncomingStream  = 0x2
)

// RPC holds any arbitrary data that is being sent over
// each transport between two nodes in the network.
type RPC struct {
	From    string // From Remote Address
	Payload []byte // Message Payload
	Steam   bool   // Stream Bool to check decoding
}
