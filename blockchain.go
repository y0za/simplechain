package main

type BlockChain struct {
	chain []Block
}

func NewBlockChain(chain ...Block) *BlockChain {
	return &BlockChain{chain}
}

func (bc *BlockChain) LatestBlock() *Block {
	if bc.chain == nil {
		return nil
	}

	l := len(bc.chain)
	if l == 0 {
		return nil
	}

	return &bc.chain[l-1]
}
