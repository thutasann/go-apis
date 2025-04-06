package p2p

import (
	"encoding/gob"
	"io"
)

// Decoder
type Decoder interface {
	Decode(io.Reader, any) error // Decode Function
}

// GOB Decoder
type GOBDecoder struct{}

// GOB Decoder
func (dec GOBDecoder) Decode(r io.Reader, v interface{}) error {
	return gob.NewDecoder(r).Decode(v)
}
