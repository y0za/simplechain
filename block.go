package main

import (
	"crypto/sha256"
	"fmt"
	"time"
)

// Block is core data structure of blockchain
type Block struct {
	Index        int    `json:"index"`
	PreviousHash string `json:"previousHash"`
	Timestamp    int64  `json:"timestamp"`
	Data         string `json:"data"`
	Hash         string `json:"hash"`
}

// NextBlock create new block from previous block and data
func NextBlock(prev Block, data string) Block {
	return nextBlockWithTimestamp(prev, data, time.Now().Unix())
}

// CheckHash verify whether hash is collect
func (b Block) CheckHash() bool {
	return b.calculateHash() == b.Hash
}

func GenesisBlock() Block {
	return Block{
		Index:        1,
		PreviousHash: "0",
		Timestamp:    1465154705,
		Data:         "my genesis block!!",
		Hash:         "816534932c2b7154836da6afc367695e6337db8a921823784c14378abed4f7d7",
	}
}

func nextBlockWithTimestamp(prev Block, data string, timestamp int64) Block {
	b := Block{
		Index:        prev.Index + 1,
		PreviousHash: prev.Hash,
		Timestamp:    timestamp,
		Data:         data,
	}
	b.Hash = b.calculateHash()
	return b
}

func (b Block) calculateHash() string {
	chank := fmt.Sprintf("%d%s%d%s", b.Index, b.PreviousHash, b.Timestamp, b.Data)
	hash := sha256.Sum256([]byte(chank))
	return fmt.Sprintf("%x", hash)
}
