package core

import "github.com/thutasann/projectx/types"

type Header struct {
	Version   uint32
	PrevBlock types.Hash
}

type Block struct{}
