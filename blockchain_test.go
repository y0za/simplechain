package main

import "testing"

func TestLatestBlock(t *testing.T) {
	var bc *BlockChain
	var expected *Block

	bc = &BlockChain{nil}
	expected = nil
	if bc.LatestBlock() != expected {
		t.Errorf("expected nil, actual not nil")
	}

	bc = &BlockChain{[]Block{}}
	expected = nil
	if bc.LatestBlock() != expected {
		t.Errorf("expected nil, actual not nil")
	}

	chain := []Block{
		Block{1, "prevhash1", 100, []byte("block1"), "hash1"},
		Block{2, "prevhash2", 200, []byte("block2"), "hash2"},
	}
	bc = NewBlockChain(chain...)
	expected = &chain[1]
	if bc.LatestBlock() != expected {
		t.Errorf("unexpected bloack")
	}
}
