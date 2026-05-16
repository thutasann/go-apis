package core

import (
	"github.com/thutasann/projectx/crypto"
	"github.com/thutasann/projectx/types"
)

// Header represents the metadata of a blockchain block.
type Header struct {
	Version       uint32
	Datahash      types.Hash
	PrevBlockHash types.Hash
	Timestamp     int64
	Height        uint32
}

// Block represents a blockchain block containing a header and a list of transactions.
type Block struct {
	*Header
	Transactions []Transaction
	Validator    crypto.PublicKey
	Signature    *crypto.Signature

	// Cached version of the header hash
	hash types.Hash
}

func (b *Block) Hash(hasher Hasher[*Block]) types.Hash {
	if b.hash.IsZero() {
		b.hash = hasher.Hash(b)
	}

	return b.hash
}
