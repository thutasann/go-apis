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
	buf := make([]byte, 1028)
	n, err := r.Read(buf)

	if err != nil {
		return err
	}

	msg.Payload = buf[:n]
	return nil
}
