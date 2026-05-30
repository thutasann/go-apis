package core

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBlockChain(t *testing.T) {
	bc, err := NewBlockChain(randomBlock(0))
	assert.Nil(t, err)
	assert.NotNil(t, bc.validator)
	assert.Equal(t, bc.Height(), uint32(0))

	fmt.Println("bc.height()", bc.Height())
}
