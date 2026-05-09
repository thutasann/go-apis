package core

import "github.com/thutasann/projectx/types"

type Header struct {
	Version   uint32
	PrevBlock types.Hash
	Timestamp uint64
	Height    uint32
	Nonce     uint64
}

type Block struct {
	Header
	Transactions []Transaction
}
