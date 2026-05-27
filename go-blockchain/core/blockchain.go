package core

type BlockChain struct {
	store     Storage
	headers   []*Header
	validator Validator
}

func (bc *BlockChain) AddBlock(b *Block) error {
	return nil
}

// [0, 1, 2, 3] => 4 len => 3 height
func (bc *BlockChain) Height() uint32 {
	return uint32(len(bc.headers) - 1)
}
