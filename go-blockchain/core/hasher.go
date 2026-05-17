package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"

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

func (BlockHasher) Hash(b *Block) types.Hash {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(b.Header); err != nil {
		panic(err)
	}

	h := sha256.Sum256(buf.Bytes())
	return types.Hash(h)
}
