package core

import (
	"crypto/sha256"

	"github.com/thutasann/projectx/types"
)

// Hasher[T] is a generic interface for computing a stable hash of a value
// of type T. Implementations should return a digest appropriate for the
// blockchain's hash type.
type Hasher[T any] interface {
	Hash(T) types.Hash
}

// BlockHasher computes a block hash from the serialized block header.
// It uses gob encoding on the header and SHA-256 to produce a deterministic
// hash suitable for block identification.
type BlockHasher struct{}

func (BlockHasher) Hash(b *Header) types.Hash {
	h := sha256.Sum256(b.Bytes())
	return types.Hash(h)
}
