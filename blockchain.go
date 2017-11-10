package main

import (
	"fmt"
	"reflect"
)

// BlockChain have multiple blocks
type BlockChain struct {
	chain []Block
}

// NewBlockChain create blockchain from some blocks
func NewBlockChain(chain ...Block) *BlockChain {
	return &BlockChain{chain}
}

// LatestBlock return last block in the blockchain
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

// CheckChain verify whether all blocks are collect and linked
func (bc BlockChain) CheckChain(genesis Block) error {
	if bc.chain == nil || len(bc.chain) == 0 {
		return fmt.Errorf("must have 1 or more blocks")
	}

	if !reflect.DeepEqual(bc.chain[0], genesis) {
		return fmt.Errorf("invalid genesis block")
	}

	if len(bc.chain) == 1 {
		return nil
	}

	for i := 1; i < len(bc.chain); i++ {
		err := checkNewBlock(bc.chain[i], &bc.chain[i-1])
		if err != nil {
			return err
		}
	}
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
