package types

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// Hash represents a 32-byte hash value.
type Hash [32]uint8

// ToSlice converts the Hash to a byte slice.
func (h Hash) ToSlice() []byte {
	b := make([]byte, 32)

	for i := range 32 {
		b[i] = h[i]
	}
	return b
}

// String returns the hexadecimal string representation of the Hash.
func (h Hash) String() string {
	return hex.EncodeToString(h.ToSlice())
}

// IsZero checks if the Hash is all zeros.
func (h Hash) IsZero() bool {
	for i := range 32 {
		if h[i] != 0 {
			return false
		}
	}
	return true
}

// HashFromBytes creates a Hash from a byte slice, panicking if the length is not 32.
func HashFromBytes(b []byte) Hash {
	if len(b) != 32 {
		msg := fmt.Sprintf("given bytes with length %d should be 32", len(b))
		panic(msg)
	}

	var value [32]uint8
	for i := 0; i < 32; i++ {
		value[i] = b[i]
	}

	return Hash(value)
}

// RandomBytes generates a slice of random bytes of the specified size.
func RandomBytes(size int) []byte {
	token := make([]byte, size)
	rand.Read(token)
	return token
}

// RandomHash generates a random 32-byte hash.
func RandomHash() Hash {
	return HashFromBytes(RandomBytes(32))
}
