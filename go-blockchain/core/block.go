package core

import (
	"io"

	"github.com/thutasann/projectx/crypto"
	"github.com/thutasann/projectx/types"
)

// Header holds the canonical metadata for a block. Fields are kept minimal
// and should be sufficient to validate and order blocks in the chain.
type Header struct {
	Version       uint32
	Datahash      types.Hash
	PrevBlockHash types.Hash
	Timestamp     int64
	Height        uint32
}

// Block represents a single block in the blockchain. It embeds a Header
// pointer and contains the block's transactions, the public key of the
// validator that proposed the block, and an optional cryptographic
// signature over the block.
type Block struct {
	*Header
	Transactions []Transaction
	Validator    crypto.PublicKey
	Signature    *crypto.Signature

	// hash caches the computed header hash for the block to avoid
	// repeated expensive hash computations.
	hash types.Hash
}

// Decode reads a Block from r using the supplied Decoder implementation.
// The decoded data is written into the receiver. Any decoding error is
// returned to the caller.
func (b *Block) Decode(r io.Reader, dec Decoder[*Block]) error {
	return dec.Decode(r, b)
}

// Encode writes the Block to w using the supplied Encoder implementation.
// Any encoding or write error is returned.
func (b *Block) Encode(w io.Writer, enc Encoder[*Block]) error {
	return enc.Encode(w, b)
}

// Hash returns the block's header hash. The value is computed lazily and
// cached in the `hash` field. If the cached value is zero, the provided
// hasher is used to compute and store the hash before returning it.
func (b *Block) Hash(hasher Hasher[*Block]) types.Hash {
	if b.hash.IsZero() {
		b.hash = hasher.Hash(b)
	}

	return b.hash
}
