package core

type Validator interface {
	ValidateBlock(*Block) error
}

type BlockValidator struct {
	bc *BlockChain
}

func NewBlockValidator(bc *BlockChain) *BlockValidator {
	return &BlockValidator{bc}
}

func (v *BlockValidator) ValidateBlock(b *Block) error {
	return nil
}
