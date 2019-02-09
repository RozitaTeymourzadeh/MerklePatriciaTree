package p1

import (
	"encoding/hex"
	"errors"
	"fmt"
	"reflect"

	"golang.org/x/crypto/sha3"
)

type stack []Node

type Flag_value struct {
	encoded_prefix []uint8
	value          string
}

type Node struct {
	node_type    int // 0: Null, 1: Branch, 2: Ext or Leaf
	branch_value [17]string
	flag_value   Flag_value
}

type MerklePatriciaTrie struct {
	db   map[string]Node
	root string
}

func HexConverter(key string) []uint8 {
	var hex_array []uint8

	for i := 0; i < len(key); i++ {
		hex_array = append(hex_array, key[i]/16)
		hex_array = append(hex_array, key[i]%16)
	}
	return hex_array
}

func Compact_encode(hex_array []uint8) []uint8 {

	var term = 0
	var result []uint8
	if len(hex_array) == 0 {
		fmt.Println("Invalid input data for Compact_encode")
		return result
	}
	if hex_array[len(hex_array)-1] == 16 {
		term = 1
	}
	if term == 1 {
		hex_array = hex_array[0 : len(hex_array)-1]
	}
	var oddlen = len(hex_array) % 2
	var flags uint8 = uint8(2*term + oddlen)
	if oddlen == 1 {
		hex_array = append([]uint8{flags}, hex_array...)
	} else {
		hex_array = append([]uint8{0}, hex_array...)
		hex_array = append([]uint8{flags}, hex_array...)
	}
	for i := 0; i < len(hex_array); i += 2 {
		result = append(result, 16*hex_array[i]+hex_array[i+1])
	}
	return result
}

func GetNodeType(node Node) uint8 {

	var nodeType uint8
	encoded_prefix := node.flag_value.encoded_prefix[0]
	nodeType = encoded_prefix / 16
	return nodeType
}

func Compact_decode(encoded_arr []uint8) []uint8 {

	var decoded_arr []uint8
	if len(encoded_arr) == 0 {
		fmt.Println("Invalid input data for Compact_decode")
		return decoded_arr
	}
	for i := 0; i < len(encoded_arr); i += 1 {
		decoded_arr = append(decoded_arr, encoded_arr[i]/16)
		decoded_arr = append(decoded_arr, encoded_arr[i]%16)
	}

	switch decoded_arr[0] {
	case 0:
		decoded_arr = decoded_arr[2:len(decoded_arr)]
	case 1:
		decoded_arr = decoded_arr[1:len(decoded_arr)]
	case 2:
		decoded_arr = decoded_arr[2:len(decoded_arr)]
	case 3:
		decoded_arr = decoded_arr[1:len(decoded_arr)]
	default:
		fmt.Println("FATAL: Invalid prefix for Compac_decoder function!")
	}
	return decoded_arr
}

func Test_compact_encode() {
	fmt.Println(reflect.DeepEqual(Compact_decode(Compact_encode([]uint8{1, 2, 3, 4, 5})), []uint8{1, 2, 3, 4, 5}))
	fmt.Println(reflect.DeepEqual(Compact_decode(Compact_encode([]uint8{0, 1, 2, 3, 4, 5})), []uint8{0, 1, 2, 3, 4, 5}))
	fmt.Println(reflect.DeepEqual(Compact_decode(Compact_encode([]uint8{0, 15, 1, 12, 11, 8, 16})), []uint8{0, 15, 1, 12, 11, 8}))
	fmt.Println(reflect.DeepEqual(Compact_decode(Compact_encode([]uint8{15, 1, 12, 11, 8, 16})), []uint8{15, 1, 12, 11, 8}))
}

func InitializeMpt() *MerklePatriciaTrie {
	db := make(map[string]Node)
	root := ""
	return &MerklePatriciaTrie{db, root}
}

func InitializeStack() stack {
	s := make(stack, 0)
	return s
}

func (s stack) Push(v Node) stack {
	return append(s, v)
}

func (s stack) Pop() (stack, Node) {
	if len(s) == 0 {
		fmt.Println("Stack is Empty!!")
		return s, Node{}
	}

	l := len(s)
	return s[:l-1], s[l-1]
}
func EqualArray(a, b []uint8) (int, []uint8){
	var j int
	var remainPath []uint8
	j = -1
	if len(a) == 0 || len(b) == 0 {
		fmt.Println("FATAL: there is no path to compare...")
		return j, remainPath
	}
	if len(a) > len(b) {
		fmt.Println("FATAL: Path not found!...")
		return j, remainPath
	}
	for i, v := range a {
		if v == b[i] {
			j++
		} else {
			break
		}
	}
	if j == len(b){
		return j , remainPath
	}
	remainPath = b[j+1:]
	return j , remainPath
}

func (mpt *MerklePatriciaTrie) Get(key string) (string, error) {
	var value string
	var errorMsg error

	currentNode := mpt.db[mpt.root]
	nodeType := GetNodeType(currentNode)
	currentPath := Compact_decode(currentNode.flag_value.encoded_prefix)
	searchPath := HexConverter(key)

	matchedIndex, remainPath := EqualArray(currentPath, searchPath)
	fmt.Println("remainPath is:",remainPath)
	fmt.Println("matchedIndex is:",matchedIndex)
	
	for len(remainPath) != 0 {
		if value == "" && errorMsg == nil && nodeType < 2{
			//Extention
			value, errorMsg, remainPath, nodeType, nextNode := mpt.FindLeafNode(currentNode, searchPath) 
			fmt.Println("----- Find Leaf Node ---------")
			fmt.Println("value is:",value)
			fmt.Println("errorMsg is:",errorMsg)
			fmt.Println("remainPath is:",remainPath)
			fmt.Println("nodeType is:",nodeType)
			fmt.Println("nextNode is:",nextNode)

			searchPath = remainPath
			currentNode = nextNode
		} 
	}
		if value == "" && errorMsg == nil && nodeType >= 2 {
			//Leaf
			value, errorMsg, remainPath := FindLeafValue (currentNode , searchPath)
			fmt.Println("----- Find Leaf Value ---------")
			fmt.Println("value is:",value)
			fmt.Println("errorMsg is:",errorMsg)
			fmt.Println("remainPath is:",remainPath)
		}
	
		return value, errorMsg
}

func FindLeafValue (node Node, searchPath []uint8) (string, error, []uint8) {
	currentPath := Compact_decode(node.flag_value.encoded_prefix)
	matchedIndex, remainPath := EqualArray(currentPath, searchPath)
		// if whole match, return Leaf Node Value
		if matchedIndex+1 == len(currentPath) {
			value := node.flag_value.value
			return value, nil, remainPath
		} else {
			// if NOT match, return Leaf Node Value
			return "", errors.New("path_not_found"),remainPath
		}
}

func (mpt *MerklePatriciaTrie) FindLeafNode(node Node, searchPath []uint8) (string, error, []uint8, uint8, Node){
	 var value string
	 var remainPath []uint8 
	 var nodeType uint8
	 var nextNode Node

	currentPath := Compact_decode(node.flag_value.encoded_prefix)
	matchedIndex, remainPath := EqualArray(currentPath, searchPath)
	// if whole match, return Branch Node Value
	if matchedIndex+1 == len(currentPath) && len(remainPath) == 0{
		nextBranchNode := mpt.db[node.flag_value.value]
		value = nextBranchNode.branch_value[16]
		return value, nil, remainPath, nodeType, nextBranchNode
	}
	// if whole match path and no remaining value 
	if matchedIndex+1 == len(currentPath) && len(remainPath)!=0 {		
		nextBranchNode := mpt.db[node.flag_value.value]
		if nextBranchNode.branch_value[remainPath[0]] != "" {
					 nextNode = mpt.db[nextBranchNode.branch_value[remainPath[0]]]
					 nodeType = GetNodeType(nextNode)
		 			// if next node is Leaf 
		 				nodeType = GetNodeType(nextNode)
						remainPath = remainPath[1:]
		} else {
		// if no match path
			return "", errors.New("path_not_found"), remainPath, nodeType, nextNode
		}	
	}
	return "", nil, remainPath, nodeType, nextNode
}

func (mpt *MerklePatriciaTrie) Delete(key string) error {
	// TODO
	return errors.New("path_not_found")
}

func (mpt *MerklePatriciaTrie) InsertRoot(key string, new_value string, s stack) {
	// first Insert
	if len(s) == 0 {
		hex_array := HexConverter(key)
		hex_array = append(hex_array, 16)
		flagValue := Flag_value{
			encoded_prefix: Compact_encode(hex_array),
			value:          new_value,
		}

		newNode := Node{
			node_type:  2,
			flag_value: flagValue,
		}

		hashedNode := newNode.Hash_node()
		mpt.db[hashedNode] = newNode
		mpt.root = hashedNode
		fmt.Println("Finish add")
		//s.Push(newNode)
		return
	}

}

func (mpt *MerklePatriciaTrie) Insert(key string, new_value string) {

}

func (node *Node) Hash_node() string {
	var str string
	switch node.node_type {
	case 0:
		str = ""
	case 1:
		str = "branch_"
		for _, v := range node.branch_value {
			str += v
		}
	case 2:
		str = node.flag_value.value
	}

	sum := sha3.Sum256([]byte(str))
	return "HashStart_" + hex.EncodeToString(sum[:]) + "_HashEnd"
}

func (mpt *MerklePatriciaTrie) GetRootNode() {
	root := mpt.db[mpt.root]
	fmt.Println("Nodetype: ", root.node_type)
	fmt.Println("Branch:", root.branch_value)
	fmt.Println("Prefix:", root.flag_value.encoded_prefix)
	fmt.Println("Value: ", root.flag_value.value)
}

//key: a value: apple
func (mpt *MerklePatriciaTrie) CreateTestMpt() error {

	mpt.db = make(map[string]Node)

	//hex_array := ConvertStringToHexArray(key)
	//hex_array = append(hex_array, 16)

	nodeCPrefix := HexConverter("1")
	nodeCPrefix = append(nodeCPrefix, 16)

	// 0: Null, 1: Branch, 2: Ext or Leaf
	flagValueNodeC := Flag_value{
			encoded_prefix: Compact_encode(nodeCPrefix),
			value:          "apple",
	}

	nodeC := Node{
			node_type:  2, //Leaf
			flag_value: flagValueNodeC,
	}
	hashNodeC := nodeC.Hash_node()
	mpt.db[hashNodeC] = nodeC

	nodeDPrefix := HexConverter("2")
	nodeDPrefix = append(nodeDPrefix, 16)
	// 0: Null, 1: Branch, 2: Ext or Leaf
	flagValueNodeD := Flag_value{
			encoded_prefix: Compact_encode(nodeDPrefix),
			value:          "banana",
	}

	nodeD := Node{
			node_type:  2, //Leaf
			flag_value: flagValueNodeD,
	}
	hashNodeD := nodeD.Hash_node()
	mpt.db[hashNodeD] = nodeD

	/////////E NODE //////////////////////////
	nodeEPrefix := HexConverter("56")
	nodeEPrefix = append(nodeEPrefix, 16)
	// 0: Null, 1: Branch, 2: Ext or Leaf
	flagValueNodeE := Flag_value{
			encoded_prefix: Compact_encode(nodeEPrefix),
			value:          "google",
	}
	nodeE := Node{
			node_type:  2, //Leaf
			flag_value: flagValueNodeE,
	}
	hashNodeE := nodeE.Hash_node()
	mpt.db[hashNodeE] = nodeE

	//////////////F NODE BRANCHE-2 /////////////////
	// 0: Null, 1: Branch, 2: Ext or Leaf
	flagValueNodeFbranch := Flag_value{
			encoded_prefix: nil,
			value:          "",
	}

	nodeFbranch := Node{
			node_type:    1, //branch
			flag_value:   flagValueNodeFbranch,
			branch_value: [17]string{"", "", "", "", "", "", "", "", "", hashNodeE, "", "", "", "", "", "", ""},
	}

	hashNodeF := nodeFbranch.Hash_node()
	mpt.db[hashNodeF] = nodeFbranch

	// 0: Null, 1: Branch, 2: Ext or Leaf
	flagValueNodeBbranch := Flag_value{
			encoded_prefix: nil,
			value:          "",
	}

	/////////G EXTENSION NODE //////////////////////////
	nodeGPrefix := HexConverter("34")
	nodeGPrefix = append(nodeGPrefix, 16)
	// 0: Null, 1: Branch, 2: Ext or Leaf
	flagValueNodeG := Flag_value{
			encoded_prefix: Compact_encode(nodeGPrefix),
			value:          hashNodeF,
	}
	nodeGExtension := Node{
			node_type:  2, //Leaf
			flag_value: flagValueNodeG,
	}
	hashNodeG := nodeGExtension.Hash_node()
	mpt.db[hashNodeG] = nodeGExtension

	nodeBbranch := Node{
			node_type:    1, //branch
			flag_value:   flagValueNodeBbranch,
			branch_value: [17]string{"", hashNodeC, hashNodeD, "", "", "", "", hashNodeG, "", "", "", "", "", "", "", "", ""},
	}

	hashNodeB := nodeBbranch.Hash_node()
	mpt.db[hashNodeB] = nodeBbranch

	nodeRootPrefix := HexConverter("6")
	flagValueRootNode := Flag_value{
			encoded_prefix: Compact_encode(nodeRootPrefix),
			value:          hashNodeB,
	}

	rootNode := Node{
			node_type:  2,
			flag_value: flagValueRootNode,
	}

	hashedNode := rootNode.Hash_node()

	mpt.db[hashedNode] = rootNode
	mpt.root = hashedNode

	fmt.Println("Mpt.db.NodeC:", mpt.db[hashNodeC])
	fmt.Println("MPt:", mpt)

	//Add another node to MPT

	//mpt.root = Compact_encode(key) //a
	return errors.New("Problem occured while creating Root Node")
}
