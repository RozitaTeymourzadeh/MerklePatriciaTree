package p2

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

/* Initial
*
* To Initialize blockChain
*
 */
func (blockChain *BlockChain) Initial() {
	blockChain.Chain = make(map[int32][]Block)
	blockChain.Length = 0
}
