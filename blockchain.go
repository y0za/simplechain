package main

import (
	"errors"
	"reflect"
)

// Blockchain have multiple blocks
type Blockchain struct {
	chain []Block
}

// NewBlockchain create blockchain from some blocks
func NewBlockchain(chain ...Block) *Blockchain {
	return &Blockchain{chain}
}

// LatestBlock return last block in the blockchain
func (bc *Blockchain) LatestBlock() *Block {
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
func (bc *Blockchain) AddBlock(b Block) error {
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

// ReplaceBlocks change blockchain data to new blocks
func (bc *Blockchain) ReplaceBlocks(blocks []Block, genesis Block) error {
	err := CheckBlocks(blocks, genesis)
	if err != nil {
		return err
	}

	if len(blocks) <= len(bc.chain) {
		return errors.New("new blocks length must be longer than current")
	}

	bc.chain = blocks
	return nil
}

// CheckBlocks verify whether all blocks are collect and linked
func CheckBlocks(blocks []Block, genesis Block) error {
	if blocks == nil || len(blocks) == 0 {
		return errors.New("must have 1 or more blocks")
	}

	if !reflect.DeepEqual(blocks[0], genesis) {
		return errors.New("invalid genesis block")
	}

	if len(blocks) == 1 {
		return nil
	}

	for i := 1; i < len(blocks); i++ {
		err := checkNewBlock(blocks[i], &blocks[i-1])
		if err != nil {
			return err
		}
	}
	return nil
}

func checkNewBlock(next Block, prev *Block) error {
	if !next.CheckHash() {
		return errors.New("invalid block hash")
	}

	if prev == nil {
		return nil
	}

	if next.Index != prev.Index+1 {
		return errors.New("unexpected block index")
	}
	if next.PreviousHash != prev.Hash {
		return errors.New("invalid block previous hash")
	}
	return nil
}
