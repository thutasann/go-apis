package p2p

import (
	"encoding/gob"
	"fmt"
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

// Go Binary Decoder
func (dec GOBDecoder) Decode(r io.Reader, rpc *RPC) error {
	return gob.NewDecoder(r).Decode(rpc)
}

// Default Decoder
func (dec DefaultDecoder) Decode(r io.Reader, rpc *RPC) error {
	buf := make([]byte, 1028)
	n, err := r.Read(buf)
	if err != nil {
		return err
	}

	fmt.Println(string(buf[:n]))

	rpc.Payload = buf[:n]

	return nil
}
