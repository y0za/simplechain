package main

import (
	"crypto/sha256"
	"fmt"
	"time"
)

// Block is core data structure of blockchain
type Block struct {
	Index        int
	PreviousHash string
	Timestamp    int64
	Data         []byte
	Hash         string
}

// NextBlock create new block from previous block and data
func NextBlock(prev Block, data []byte) Block {
	b := Block{
		Index:        prev.Index + 1,
		PreviousHash: prev.Hash,
		Timestamp:    time.Now().Unix(),
		Data:         data,
	}
	b.Hash = b.calculateHash()
	return b
}

// CheckHash verify whether hash is collect
func (b Block) CheckHash() bool {
	return b.calculateHash() == b.Hash
}

func (b Block) calculateHash() string {
	chank := fmt.Sprintf("%d%s%d%s", b.Index, b.PreviousHash, b.Timestamp, b.Data)
	hash := sha256.Sum256([]byte(chank))
	return fmt.Sprintf("%x", hash)
}
