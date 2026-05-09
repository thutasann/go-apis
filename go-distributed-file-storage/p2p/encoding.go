package p2p

import (
	"encoding/gob"
	"io"
)

// Decoder Interface
type Decoder interface {
	Decode(io.Reader, *RPC) error // Decode Function
}

// Go Binary Decoder
type GOBDecoder struct{}

// Default Decoder
type DefaultDecoder struct{}

// Go Binary Decoder implements Decoder interface
func (dec GOBDecoder) Decode(r io.Reader, msg *RPC) error {
	return gob.NewDecoder(r).Decode(msg)
}

// Default Decoder implements Decoder interface
func (dec DefaultDecoder) Decode(r io.Reader, msg *RPC) error {
	peekBuf := make([]byte, 1)
	if _, err := r.Read(peekBuf); err != nil {
		return nil
	}

	// in case of a stream we are not decoding what is being sent over the network,
	// we are just setting stream true so we can handle that in our logic
	stream := peekBuf[0] == IncomingStream
	if stream {
		msg.Steam = true
		return nil
	}

	buf := make([]byte, 1028)
	n, err := r.Read(buf)

	if err != nil {
		return err
	}

	msg.Payload = buf[:n]
	return nil
}
