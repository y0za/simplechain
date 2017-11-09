package main

import (
	"reflect"
	"testing"
)

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

func TestAddBlock(t *testing.T) {
	tests := []struct {
		bc        *BlockChain
		b         Block
		expected  *BlockChain
		expectErr bool
	}{
		{
			NewBlockChain(Block{1, "prevhash1", 100, []byte("block1"), "hash1"}),
			Block{2, "prevhash2", 200, []byte("block2"), "hash2"},
			NewBlockChain(Block{1, "prevhash1", 100, []byte("block1"), "hash1"}),
			true,
		},
		{
			NewBlockChain(Block{0, "prevhash1", 100, []byte("block1"), "hash1"}),
			nextBlockWithTimestamp(Block{1, "prevhash1", 100, []byte("block1"), "hash1"}, []byte("block2"), 200),
			NewBlockChain(Block{0, "prevhash1", 100, []byte("block1"), "hash1"}),
			true,
		},
		{
			NewBlockChain(Block{0, "prevhash1", 100, []byte("block1"), "invalid hash"}),
			nextBlockWithTimestamp(Block{1, "prevhash1", 100, []byte("block1"), "hash1"}, []byte("block2"), 200),
			NewBlockChain(Block{0, "prevhash1", 100, []byte("block1"), "invalid hash"}),
			true,
		},
		{
			&BlockChain{nil},
			nextBlockWithTimestamp(Block{1, "prevhash1", 100, []byte("block1"), "hash1"}, []byte("block2"), 200),
			NewBlockChain(
				nextBlockWithTimestamp(Block{1, "prevhash1", 100, []byte("block1"), "hash1"}, []byte("block2"), 200),
			),
			false,
		},
		{
			&BlockChain{[]Block{}},
			nextBlockWithTimestamp(Block{1, "prevhash1", 100, []byte("block1"), "hash1"}, []byte("block2"), 200),
			NewBlockChain(
				nextBlockWithTimestamp(Block{1, "prevhash1", 100, []byte("block1"), "hash1"}, []byte("block2"), 200),
			),
			false,
		},
		{
			NewBlockChain(Block{1, "prevhash1", 100, []byte("block1"), "hash1"}),
			nextBlockWithTimestamp(Block{1, "prevhash1", 100, []byte("block1"), "hash1"}, []byte("block2"), 200),
			NewBlockChain(
				Block{1, "prevhash1", 100, []byte("block1"), "hash1"},
				nextBlockWithTimestamp(Block{1, "prevhash1", 100, []byte("block1"), "hash1"}, []byte("block2"), 200),
			),
			false,
		},
	}

	for i, tt := range tests {
		err := tt.bc.AddBlock(tt.b)
		if tt.expectErr && err == nil {
			t.Errorf("case %d expected error, actual nil", i)
		}
		if !tt.expectErr && err != nil {
			t.Errorf("case %d expected no error, actual error '%s'", i, err)
		}
		if !reflect.DeepEqual(tt.expected, tt.bc) {
			t.Errorf("case %d unexpected result chain", i)
		}
	}
}

func TestCheckNewBlock(t *testing.T) {
	tests := []struct {
		next      Block
		prev      *Block
		expectErr bool
	}{
		{
			Block{2, "prevhash2", 200, []byte("block2"), "hash2"},
			&Block{1, "prevhash1", 100, []byte("block1"), "hash1"},
			true,
		},
		{
			NextBlock(Block{1, "prevhash1", 100, []byte("block1"), "hash1"}, []byte("block2")),
			&Block{0, "prevhash1", 100, []byte("block1"), "hash1"},
			true,
		},
		{
			NextBlock(Block{1, "prevhash1", 100, []byte("block1"), "hash1"}, []byte("block2")),
			&Block{1, "prevhash1", 100, []byte("block1"), "invalid hash"},
			true,
		},
		{
			NextBlock(Block{1, "prevhash1", 100, []byte("block1"), "hash1"}, []byte("block2")),
			nil,
			false,
		},
		{
			NextBlock(Block{1, "prevhash1", 100, []byte("block1"), "hash1"}, []byte("block2")),
			&Block{1, "prevhash1", 100, []byte("block1"), "hash1"},
			false,
		},
	}

	for _, tt := range tests {
		err := checkNewBlock(tt.next, tt.prev)
		if tt.expectErr && err == nil {
			t.Error("expected error, actual nil")
		}
		if !tt.expectErr && err != nil {
			t.Errorf("expected no error, actual error '%s'", err)
		}
	}
}
