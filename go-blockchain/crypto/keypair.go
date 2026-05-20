package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"math/big"

	"github.com/thutasann/projectx/types"
)

type PrivateKey struct {
	key *ecdsa.PrivateKey
}

// PublicKey returns the public key corresponding to this private key.
func (k PrivateKey) PublicKey() PublicKey {
	return PublicKey{
		key: &k.key.PublicKey,
	}
}

// Sign creates an ECDSA signature for the provided data using the private key.
// The data must be hashed before signing when used in higher-level protocols.
func (k PrivateKey) Sign(data []byte) (*Signature, error) {
	r, s, err := ecdsa.Sign(rand.Reader, k.key, data)
	if err != nil {
		return nil, err
	}

	return &Signature{
		r: r,
		s: s,
	}, nil
}

// GeneratePrivateKey creates a new ECDSA private key on the P-256 curve.
func GeneratePrivateKey() PrivateKey {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}

	return PrivateKey{
		key: key,
	}
}

type PublicKey struct {
	key *ecdsa.PublicKey
}

// ToSlice serializes the public key in compressed elliptic curve form.
func (k PublicKey) ToSlice() []byte {
	return elliptic.MarshalCompressed(k.key, k.key.X, k.key.Y)
}

// Address derives a blockchain address from the public key by hashing the
// compressed key and taking the last 20 bytes.
func (k PublicKey) Address() types.Address {
	h := sha256.Sum256(k.ToSlice())
	return types.NewAddressFromBytes(h[len(h)-20:])
}

type Signature struct {
	r, s *big.Int
}

// Verify checks whether the signature is valid for the given public key and data.
// The `data` parameter must be the hashed message that was originally signed.
// Callers are responsible for hashing (for example, SHA-256) before verification
// when using higher-level protocols. Returns true if the signature is valid.
func (sig Signature) Verify(pubKey PublicKey, data []byte) bool {
	return ecdsa.Verify(pubKey.key, data, sig.r, sig.s)
}
