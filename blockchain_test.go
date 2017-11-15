package main

import (
	"reflect"
	"testing"
)

func TestLatestBlock(t *testing.T) {
	var bc *Blockchain
	var expected *Block

	bc = &Blockchain{nil}
	expected = nil
	if bc.LatestBlock() != expected {
		t.Errorf("expected nil, actual not nil")
	}

	bc = &Blockchain{[]Block{}}
	expected = nil
	if bc.LatestBlock() != expected {
		t.Errorf("expected nil, actual not nil")
	}

	chain := []Block{
		Block{1, "prevhash1", 100, []byte("block1"), "hash1"},
		Block{2, "prevhash2", 200, []byte("block2"), "hash2"},
	}
	bc = NewBlockchain(chain...)
	expected = &chain[1]
	if bc.LatestBlock() != expected {
		t.Errorf("unexpected bloack")
	}
}

func TestAddBlock(t *testing.T) {
	tests := []struct {
		bc        *Blockchain
		b         Block
		expected  *Blockchain
		expectErr bool
	}{
		{
			NewBlockchain(Block{1, "prevhash1", 100, []byte("block1"), "hash1"}),
			Block{2, "prevhash2", 200, []byte("block2"), "hash2"},
			NewBlockchain(Block{1, "prevhash1", 100, []byte("block1"), "hash1"}),
			true,
		},
		{
			NewBlockchain(Block{0, "prevhash1", 100, []byte("block1"), "hash1"}),
			nextBlockWithTimestamp(Block{1, "prevhash1", 100, []byte("block1"), "hash1"}, []byte("block2"), 200),
			NewBlockchain(Block{0, "prevhash1", 100, []byte("block1"), "hash1"}),
			true,
		},
		{
			NewBlockchain(Block{0, "prevhash1", 100, []byte("block1"), "invalid hash"}),
			nextBlockWithTimestamp(Block{1, "prevhash1", 100, []byte("block1"), "hash1"}, []byte("block2"), 200),
			NewBlockchain(Block{0, "prevhash1", 100, []byte("block1"), "invalid hash"}),
			true,
		},
		{
			&Blockchain{nil},
			nextBlockWithTimestamp(Block{1, "prevhash1", 100, []byte("block1"), "hash1"}, []byte("block2"), 200),
			NewBlockchain(
				nextBlockWithTimestamp(Block{1, "prevhash1", 100, []byte("block1"), "hash1"}, []byte("block2"), 200),
			),
			false,
		},
		{
			&Blockchain{[]Block{}},
			nextBlockWithTimestamp(Block{1, "prevhash1", 100, []byte("block1"), "hash1"}, []byte("block2"), 200),
			NewBlockchain(
				nextBlockWithTimestamp(Block{1, "prevhash1", 100, []byte("block1"), "hash1"}, []byte("block2"), 200),
			),
			false,
		},
		{
			NewBlockchain(Block{1, "prevhash1", 100, []byte("block1"), "hash1"}),
			nextBlockWithTimestamp(Block{1, "prevhash1", 100, []byte("block1"), "hash1"}, []byte("block2"), 200),
			NewBlockchain(
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

func TestCheckBlocks(t *testing.T) {
	tests := []struct {
		blocks    []Block
		genesis   Block
		expectErr bool
	}{
		{
			nil,
			Block{1, "prevhash1", 100, []byte("block1"), "hash1"},
			true,
		},
		{
			[]Block{},
			Block{1, "prevhash1", 100, []byte("block1"), "hash1"},
			true,
		},
		{
			[]Block{Block{1, "prevhash1", 100, []byte("not genesis"), "hash1"}},
			Block{1, "prevhash1", 100, []byte("block1"), "hash1"},
			true,
		},
		{
			[]Block{
				Block{1, "prevhash1", 100, []byte("block1"), "hash1"},
				Block{2, "hash1", 200, []byte("block1"), "invalid hash"},
			},
			Block{1, "prevhash1", 100, []byte("block1"), "hash1"},
			true,
		},
		{
			[]Block{Block{1, "prevhash1", 100, []byte("block1"), "hash1"}},
			Block{1, "prevhash1", 100, []byte("block1"), "hash1"},
			false,
		},
		{
			[]Block{
				Block{1, "prevhash1", 100, []byte("block1"), "hash1"},
				nextBlockWithTimestamp(Block{1, "prevhash1", 100, []byte("block1"), "hash1"}, []byte("block2"), 200),
			},
			Block{1, "prevhash1", 100, []byte("block1"), "hash1"},
			false,
		},
	}

	for i, tt := range tests {
		err := CheckBlocks(tt.blocks, tt.genesis)
		if tt.expectErr && err == nil {
			t.Errorf("case %d expected error, actual nil", i)
		}
		if !tt.expectErr && err != nil {
			t.Errorf("case %d expected no error, actual error '%s'", i, err)
		}
	}
}

func TestReplaceBlocks(t *testing.T) {
	tests := []struct {
		bc        *Blockchain
		blocks    []Block
		genesis   Block
		expectErr bool
	}{
		{
			NewBlockchain(Block{1, "prevhash1", 100, []byte("block1"), "hash1"}),
			[]Block{
				Block{1, "prevhash1", 100, []byte("not genesis"), "hash1"},
				nextBlockWithTimestamp(Block{1, "prevhash1", 100, []byte("block1"), "hash1"}, []byte("block2"), 200),
			},
			Block{1, "prevhash1", 100, []byte("block1"), "hash1"},
			true,
		},
		{
			NewBlockchain(Block{1, "prevhash1", 100, []byte("block1"), "hash1"}),
			[]Block{
				Block{1, "prevhash1", 100, []byte("block1"), "hash1"},
			},
			Block{1, "prevhash1", 100, []byte("block1"), "hash1"},
			true,
		},
		{
			NewBlockchain(Block{1, "prevhash1", 100, []byte("block1"), "hash1"}),
			[]Block{
				Block{1, "prevhash1", 100, []byte("block1"), "hash1"},
				nextBlockWithTimestamp(Block{1, "prevhash1", 100, []byte("block1"), "hash1"}, []byte("block2"), 200),
			},
			Block{1, "prevhash1", 100, []byte("block1"), "hash1"},
			false,
		},
	}

	for i, tt := range tests {
		err := tt.bc.ReplaceBlocks(tt.blocks, tt.genesis)
		if tt.expectErr && err == nil {
			t.Errorf("case %d expected error, actual nil", i)
			return
		}
		if !tt.expectErr && err != nil {
			t.Errorf("case %d expected no error, actual error '%s'", i, err)
			return
		}
		if !tt.expectErr && !reflect.DeepEqual(tt.bc.chain, tt.blocks) {
			t.Errorf("case %d expected blockchain is replaced with new blocks, actual not", i)
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

	for i, tt := range tests {
		err := checkNewBlock(tt.next, tt.prev)
		if tt.expectErr && err == nil {
			t.Errorf("case %d expected error, actual nil", i)
		}
		if !tt.expectErr && err != nil {
			t.Errorf("case %d expected no error, actual error '%s'", i, err)
		}
	}
}
