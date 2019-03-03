package p2

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"golang.org/x/crypto/sha3"
)

/*-------------------------STRUCT---------------------------------------------------*/
/* Struct data structure for variables
/*-------------------------STRUCT---------------------------------------------------*/

/* Block Header struct
*
* To Define block header using Block header struct
*
 */
type BlockHeader struct {
	Size       int32
	ParentHash string
	Height     int32
	Timestamp  int64
	Hash       string
}

/* Block struct
*
* Block struct
 */
type Block struct {
	Header BlockHeader
	Value  MerklePatriciaTrie
}

/*-------------------------INITIALIZATION---------------------------------------------------*/
/* Data initialization
/*-------------------------INITIALIZATION---------------------------------------------------*/

/* Initial
*
* To initialize MPT height and parentHash
* @input: value MerklePatriciaTrie, height int32, parentHash string
* @output: nill
*
 */
func (block *Block) Initial(height int32, parentHash string, value MerklePatriciaTrie) {
	block.Header.Height = height
	block.Header.Timestamp = time.Now().Unix()
	block.Header.ParentHash = parentHash
	block.Value = value
}

/*-------------------------JSON HELPER---------------------------------------------------*/
/* JSON feature
/*-------------------------JSON HELPER---------------------------------------------------*/

/* UnmarshalJSON
* Interitted from golang library
* To encodes a block instance into a JSON format string
* @input: an instanse of block
* @output: a string of JSON format
*
 */
func (block *Block) UnmarshalJSON(data []byte) error {
	JsonBlock := SymmetricJsonBlock{}
	err := json.Unmarshal(data, &JsonBlock)
	if err != nil {
		return err
	}
	block.Header.Height = JsonBlock.Height
	block.Header.Timestamp = JsonBlock.Timestamp
	block.Header.Hash = JsonBlock.Hash
	block.Header.ParentHash = JsonBlock.ParentHash
	block.Header.Size = JsonBlock.Size
	mpt := MerklePatriciaTrie{}
	mpt.Initial()
	for k, v := range JsonBlock.Mpt {
		mpt.Insert(k, v)
	}
	block.Value = mpt
	return nil
}

/* EncodeToJSON
* Interitted from golang library
* To encodes a block instance into a JSON format string
* @input: an instanse of block
* @output: a string of JSON format
*
 */
func (block *Block) EncodeToJSON() (string, error) {
	jsonBytes, err := json.Marshal(block)
	return string(jsonBytes), err
}

/* DecodeFromJSON
*
* To take a string that represents the JSON value of a block as an input, and decodes the input string back to a block instance.
* @input:  a string of JSON format
* @output: an instanse of block
*
 */
func (block *Block) DecodeFromJSON(jsonString string) error {
	return json.Unmarshal([]byte(jsonString), block)
}

/* MarshalJSON
*
* To hash MPT with the SHA3-256 encoded value of this string and update MPT value upon the
* insertion
* @input:  block
* @output: updated block
*
 */
func (block *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(SymmetricJsonBlock{
		Size:       block.Header.Size,
		ParentHash: block.Header.ParentHash,
		Height:     block.Header.Height,
		Timestamp:  block.Header.Timestamp,
		Hash:       block.Header.Hash,
		Mpt:        block.Value.LeafList(),
	})
}

/* SymmetricJsonBlock
*
* Symmetric Json struct for Block struct
*
 */
type SymmetricJsonBlock struct {
	Size       int32             `json:"size"`
	ParentHash string            `json:"parentHash"`
	Height     int32             `json:"height"`
	Timestamp  int64             `json:"timeStamp"`
	Hash       string            `json:"hash"`
	Mpt        map[string]string `json:"mpt"`
}

/*-------------------------MASTER---------------------------------------------------*/
/* Function to modify Block
/*-------------------------MASTER---------------------------------------------------*/

/* Insert
*
* To insert key/value pair into the main MPT as block and
* update the mpt size and hash value. take a string that represents the JSON value of a block as
* @input:  key string, value string
* @output: update block
*
 */
func (block *Block) Insert(key string, value string) {
	block.Value.Insert(key, value)
	block.UpdateMpt()
}

/* UpdateMpt
*
* To hash MPT with the SHA3-256 encoded value of this string and update MPT value upon the
* insertion
* @input:  block
* @output: updated block
*
 */
func (block *Block) UpdateMpt() {
	block.Header.Size = int32(len([]byte(fmt.Sprintf("%v", block.Value))))
	hashConverter := sha3.New256()
	hashStr := string(block.Header.Height) + string(block.Header.Timestamp) + block.Header.ParentHash + block.Value.root + string(block.Header.Size)
	block.Header.Hash = hex.EncodeToString(hashConverter.Sum([]byte(hashStr)))
}
