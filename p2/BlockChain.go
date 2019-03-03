package p2

import "encoding/json"

/*-------------------------STRUCT---------------------------------------------------*/
/* Struct data structure for variables
/*-------------------------STRUCT---------------------------------------------------*/

/* BlockChain struct
*
* To Define blockChain variables
*
 */
type BlockChain struct {
	Chain  map[int32][]Block `json:"chain"`
	Length int32             `json:"length"`
}

/*-------------------------INITIALIZATION---------------------------------------------------*/
/* Initialize blockChain
/*-------------------------INITIALIZATION---------------------------------------------------*/

/* Initial
*
* To Initialize blockChain
*
 */
func (blockChain *BlockChain) Initial() {
	blockChain.Chain = make(map[int32][]Block)
	blockChain.Length = 0
}

/*-------------------------MASTER---------------------------------------------------*/
/* Main function
/*-------------------------MASTER---------------------------------------------------*/

/* Get
*
* To return blocks in chain with certain height
* @input: height int32
* @output: blockChain.Chain[height]
 */
func (blockChain *BlockChain) Get(height int32) []Block {
	return blockChain.Chain[height]
}

/* Insert
*
* To insert block into blockchain
*
 */
func (blockChain *BlockChain) Insert(block Block) {
	if block.Header.Height > blockChain.Length {
		blockChain.Length = block.Header.Height
	}
	heightBlocks := blockChain.Chain[block.Header.Height]
	if heightBlocks == nil { // return empty block if heght is zero
		heightBlocks = []Block{}
	}
	for _, heightBlock := range heightBlocks { // find simmilar hash in blockchain
		if heightBlock.Header.Hash == block.Header.Hash {
			return
		}
	}
	// append to blockChain
	blockChain.Chain[block.Header.Height] = append(heightBlocks, block)
}

/*-------------------------JSON HELPER---------------------------------------------------*/
/* Serialize and decerialization
/*-------------------------JSON HELPER---------------------------------------------------*/
/* UnmarshalJSON
* Interitted from golang library
* To decerialize blockChain as Json type
*
 */
func (blockChain *BlockChain) UnmarshalJSON(data []byte) error {
	blocks := make([]Block, 0)
	err := json.Unmarshal(data, &blocks)
	if err != nil {
		return err
	}
	blockChain.Initial()
	for _, block := range blocks {
		blockChain.Insert(block) // TODO
	}
	return nil //return nil if no block
}

/* MarshalJSON
* Interitted from golang library
* To serialize  blockChain as Json type
*
 */
func (blockChain *BlockChain) MarshalJSON() ([]byte, error) {
	blocks := make([]Block, 0)
	for _, v := range blockChain.Chain {
		blocks = append(blocks, v...)
	}
	return json.Marshal(blocks)
}
