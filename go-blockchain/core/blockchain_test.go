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

func TestAddBlock(t *testing.T) {
	bc := newBlockChainWithGenesif(t)

	lenBlocks := 1000
	for i := range lenBlocks {
		block := randomBlockWithSignature(t, uint32(i+1))
		assert.Nil(t, bc.AddBlock(block))
	}

	assert.Equal(t, bc.Height(), uint32(lenBlocks))
	assert.Equal(t, len(bc.headers), lenBlocks+1)
	assert.NotNil(t, bc.AddBlock(randomBlock(89)))
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

func TestAddBlockToHigh(t *testing.T) {
	bc := newBlockChainWithGenesif(t)
	assert.NotNil(t, bc.AddBlock(randomBlockWithSignature(t, 3)))
}

func TestGetHeader(t *testing.T) {
	bc := newBlockChainWithGenesif(t)
	lenBlocks := 1000

	for i := 0; i < lenBlocks; i++ {
		block := randomBlockWithSignature(t, uint32(i+1))
		assert.Nil(t, bc.AddBlock(block))
		header, err := bc.GetHeader(block.Height)
		assert.Nil(t, err)
		assert.Equal(t, header, block.Header)
	}
}
