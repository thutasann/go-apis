package core

import "github.com/thutasann/projectx/crypto"

type Transaction struct {
	Data      []byte
	PublicKey crypto.PublicKey
	Signature *crypto.Signature
}
