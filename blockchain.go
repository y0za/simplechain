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
