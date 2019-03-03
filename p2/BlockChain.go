package p2


/*-------------------------STRUCT---------------------------------------------------*/
/* Struct data structure for variables
/*-------------------------STRUCT---------------------------------------------------*/
type BlockChain struct {
	Chain  map[int32][]Block `json:"chain"`
	Length int32             `json:"length"`
}
