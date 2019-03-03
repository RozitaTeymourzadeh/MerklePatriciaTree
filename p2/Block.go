package p2

import (
	"encoding/json"
	"time"
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
	Height     int32  `json:"height"`
	Timestamp  int64  `json:"timestamp"`
	Hash       string `json:"hash"`
	ParentHash string `json:"parentHash"`
	Size       int32  `json:"size"`
}

/* Block struct
*
* Block struct
 */
type Block struct {
	Header BlockHeader        `json:"header"`
	Value  MerklePatriciaTrie `json:"value"`
}

/* DataInitialization
*
* To initialize MPT height and parentHash
* @input: value MerklePatriciaTrie, height int32, parentHash string
* @output: nill
*
 */
func (block *Block) Initial(value MerklePatriciaTrie, height int32, parentHash string) {
	block.Header.Height = height
	block.Header.Timestamp = time.Now().Unix()
	block.Header.ParentHash = parentHash
	block.Value = value
}

/* EncodeToJson
*
* To encodes a block instance into a JSON format string
* @input: an instanse of block
* @output: a string of JSON format
*
 */
func (block *Block) EncodeToJson() (string, error) {
	jsonBytes, err := json.Marshal(block)
	return string(jsonBytes), err
}

/* DecodeFromJson
*
* To take a string that represents the JSON value of a block as an input, and decodes the input string back to a block instance.
* @input:  a string of JSON format
* @output: an instanse of block
*
 */
func (block *Block) DecodeFromJson(jsonString string) error {
	return json.Unmarshal([]byte(jsonString), block)
}
