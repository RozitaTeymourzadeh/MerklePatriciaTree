package p2

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
	"golang.org/x/crypto/sha3"
)

type BlockHeader struct {
	Height     int32  `json:"height"`
	Timestamp  int64  `json:"timestamp"`
	Hash       string `json:"hash"`
	ParentHash string `json:"parentHash"`
	Size       int32  `json:"size"`
}

