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

/* UnmarshalJSON
*
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
		blockChain.Insert(block)// TODO
	}

	return nil
}
