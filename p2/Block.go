package p2

import (
	"time"
)

/* Define block header */
type BlockHeader struct {
	Height     int32  `json:"height"`
	Timestamp  int64  `json:"timestamp"`
	Hash       string `json:"hash"`
	ParentHash string `json:"parentHash"`
	Size       int32  `json:"size"`
}

/* Define block */
type Block struct {
	Header BlockHeader        `json:"header"`
	Value  MerklePatriciaTrie `json:"value"`
}

/* Initialization */
func (block *Block) DataInitialization(value MerklePatriciaTrie, height int32, parentHash string) {
	block.Header.Height = height
	block.Header.Timestamp = time.Now().Unix()
	block.Header.ParentHash = parentHash
	block.Value = value
}
