package main

import "fmt"

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

// AddBlock check block and append it to the end of the chain
func (bc *BlockChain) AddBlock(b Block) error {
	if bc.chain == nil {
		bc.chain = []Block{}
	}

	err := checkNewBlock(b, bc.LatestBlock())
	if err != nil {
		return err
	}

	bc.chain = append(bc.chain, b)
	return nil
}

func checkNewBlock(next Block, prev *Block) error {
	if !next.CheckHash() {
		return fmt.Errorf("invalid block hash")
	}

	if prev == nil {
		return nil
	}

	if next.Index != prev.Index+1 {
		return fmt.Errorf("unexpected block index")
	}
	if next.PreviousHash != prev.Hash {
		return fmt.Errorf("invalid block previous hash")
	}
	return nil
}
