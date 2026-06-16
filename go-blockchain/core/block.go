package core

import (
	"bytes"
	"encoding/gob"
	"fmt"
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

func (h *Header) Bytes() []byte {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	enc.Encode(h)

	return buf.Bytes()
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

// Initialize a new Block
func NewBlock(h *Header, txx []Transaction) *Block {
	return &Block{
		Header:       h,
		Transactions: txx,
	}
}

// Sign signs the block's header with the provided private key. It computes a
// cryptographic signature over the block's header data and stores both the
// signature and the validator's public key in the block. Returns an error if
// the signing operation fails.
func (b *Block) Sign(privKey crypto.PrivateKey) error {
	sig, err := privKey.Sign(b.Header.Bytes())
	if err != nil {
		return err
	}

	b.Validator = privKey.PublicKey()
	b.Signature = sig
	return nil
}

func (b *Block) Verify() error {
	if b.Signature == nil {
		return fmt.Errorf("block has no signature")
	}

	if !b.Signature.Verify(b.Validator, b.Header.Bytes()) {
		return fmt.Errorf("block has invalid signature")
	}

	return nil
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
func (b *Block) Hash(hasher Hasher[*Header]) types.Hash {
	if b.hash.IsZero() {
		b.hash = hasher.Hash(b.Header)
	}

	return b.hash
}
