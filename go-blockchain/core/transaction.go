package core

import (
	"fmt"

	"github.com/thutasann/projectx/crypto"
)

// Transaction represents a signed payload that can be verified by peers.
// It contains the raw data, the public key of the signer, and the signature.
type Transaction struct {
	Data      []byte
	PublicKey crypto.PublicKey
	Signature *crypto.Signature
}

// Sign signs the transaction data with the provided private key.
// It stores the signer's public key and the generated signature on the transaction.
func (tx *Transaction) Sign(privKey crypto.PrivateKey) error {
	s, err := privKey.Sign(tx.Data)
	if err != nil {
		return err
	}

	tx.PublicKey = privKey.PublicKey()
	tx.Signature = s

	return nil
}

func (tx *Transaction) Verify() error {
	if tx.Signature == nil {
		return fmt.Errorf("transaction has no signature")
	}

	if !tx.Signature.Verify(tx.PublicKey, tx.Data) {
		return fmt.Errorf("invalid transaction signature")
	}

	return nil
}
