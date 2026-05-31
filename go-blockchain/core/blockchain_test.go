package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func newBlockChainWithGenesif(t *testing.T) *BlockChain {
	bc, err := NewBlockChain(randomBlock(0))
	assert.Nil(t, err)
	return bc
}

func TestNewBlockChain(t *testing.T) {
	bc := newBlockChainWithGenesif(t)
	assert.NotNil(t, bc.validator)
	assert.Equal(t, bc.Height(), uint32(0))
}

func TestHasBlock(t *testing.T) {
	bc := newBlockChainWithGenesif(t)
	assert.True(t, bc.HasBlock(0))
}
