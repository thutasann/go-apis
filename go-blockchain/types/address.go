package types

import (
	"encoding/hex"
	"fmt"
)

type Address [20]uint8

// ToSlice returns the address as a byte slice.
func (a Address) ToSlice() []byte {
	b := make([]byte, 20)

	for i := range 20 {
		b[i] = a[i]
	}

	return b
}

// NewAddressFromBytes constructs an Address from a 20-byte slice.
// It panics if the provided byte slice is not exactly 20 bytes long.
func NewAddressFromBytes(b []byte) Address {
	if len(b) != 20 {
		msg := fmt.Sprintf("given bytes with length %d should be 20", len(b))
		panic(msg)
	}

	var value [20]uint8
	for i := range 20 {
		value[i] = b[i]
	}

	return Address(value)
}

// String returns the hexadecimal string representation of the address.
func (a Address) String() string {
	return hex.EncodeToString(a.ToSlice())
}
