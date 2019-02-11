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
func EqualArray(a, b []uint8) (int, []uint8) {
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
	if j == len(b) {
		return j, remainPath
	}
	remainPath = b[j+1:]
	return j, remainPath
}

func (mpt *MerklePatriciaTrie) Get(key string) (string, error) {
	var DEBUG int
	DEBUG = 0
	var value string
	var errorMsg error
	var nodeType uint8 
	nodeType = 100
 if key == "" {
	 fmt.Println("FATAl: Invalid KEY ...")
	 return "", errors.New("path_not_found")
 }
	currentNode := mpt.db[mpt.root]
	if currentNode.node_type == 2{
		nodeType = GetNodeType(currentNode)
	}
	searchPath := HexConverter(key)
	for len(searchPath) != 0 && value == "" && errorMsg == nil {
		// if node_type is NULL, Branch, Leaf/Ext
		if nodeType < 2 {
			//Extention
			valueFln, errorMsgFln, remainPath, nodeTypeFln, nextNode := mpt.FindLeafNode(currentNode, searchPath)
			if  DEBUG == 1{
			fmt.Println("----- Find Leaf Node ---------")
			fmt.Println("value is:", valueFln)
			fmt.Println("errorMsg is:", errorMsgFln)
			fmt.Println("remainPath is:", remainPath)
			fmt.Println("nodeType is:", nodeTypeFln)
			fmt.Println("nextNode is:", nextNode)
			}
			searchPath = remainPath
			currentNode = nextNode
			errorMsg = errorMsgFln
			nodeType = nodeTypeFln
			value = valueFln
		} else {
			//Leaf
			valueFlv, errorMsgFlv, remainPath := FindLeafValue(currentNode, searchPath)
			if  DEBUG == 1{
			fmt.Println("----- Find Leaf Value ---------")
			fmt.Println("value is:", valueFlv)
			fmt.Println("errorMsg is:", errorMsgFlv)
			fmt.Println("remainPath is:", remainPath)
			}
			value = valueFlv
			errorMsg = errorMsgFlv
		}
	} 
	return value, errorMsg
}

func FindLeafValue(node Node, searchPath []uint8) (string, error, []uint8) {
	currentPath := Compact_decode(node.flag_value.encoded_prefix)
	matchedIndex, remainPath := EqualArray(currentPath, searchPath)
	// if whole match, return Leaf Node Value
	if matchedIndex+1 == len(currentPath) {
		value := node.flag_value.value
		return value, nil, remainPath
	} else {
		// if NOT match, return Leaf Node Value
		return "", errors.New("path_not_found"), remainPath
	}
}

func (mpt *MerklePatriciaTrie) FindLeafNode(node Node, searchPath []uint8) (string, error, []uint8, uint8, Node) {
	var value string
	var remainPath []uint8
	var nodeType uint8
	var nextNode Node
	var currentPath []uint8
	var index uint8
	
	if node.node_type == 2 {
		// Extension or Leaf 
		currentPath = Compact_decode(node.flag_value.encoded_prefix)
	} else if node.node_type == 1 {
		// Branch 
		branchCurrentPath := node.branch_value
		for  index = 0; index < uint8(len(branchCurrentPath)) ; index++ {
			if branchCurrentPath[index] != "" {
				currentPath = append(currentPath, index)
			}
		}
	} else {
		// Null 
		fmt.Println("Null Branch!!")
		return "", errors.New("path_not_found"), remainPath, nodeType, nextNode
	}
	matchedIndex, remainPath := EqualArray(currentPath, searchPath)
	// if whole match, return Branch Node Value if current Node is Ext/Leaf
	if matchedIndex+1 == len(currentPath) && len(remainPath) == 0 && node.node_type == 2 {
		nextBranchNode := mpt.db[node.flag_value.value]
		value = nextBranchNode.branch_value[16]
		return value, nil, remainPath, nodeType, nextBranchNode
	}
	// if whole match, if current Node is Branch
	if matchedIndex+1 == len(currentPath) && len(remainPath) == 0 && node.node_type == 1 {
		nextNode := mpt.db[node.branch_value[currentPath[0]]]
		if nextNode.node_type == 2{
			value = nextNode.flag_value.value
		}
		return value, nil, remainPath, nodeType, nextNode
	}
	// if whole match path and remaining value
	if matchedIndex+1 == len(currentPath) && len(remainPath) != 0 {
		nextBranchNode := mpt.db[node.flag_value.value]
		if nextBranchNode.branch_value[remainPath[0]] != "" {
			nextNode = mpt.db[nextBranchNode.branch_value[remainPath[0]]]
			if nextNode.node_type == 1{
				// Node is Branch
				 remainPath = remainPath[1:]
			} else if nextNode.node_type == 2 {
				// Leaf or Extension
				nodeType = GetNodeType(nextNode)
				if nodeType  < 2 {
					//Extension
					remainPath = remainPath[1:]
				} else {
					//Leaf
					remainPath = remainPath[1:]
				}
			}
		} else {
			// if no match path
			return "", errors.New("path_not_found"), remainPath, nodeType, nextNode
		}
	}
	return "", nil , remainPath, nodeType, nextNode
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
func Test_compact_encode() {
	fmt.Println(reflect.DeepEqual(Compact_decode(Compact_encode([]uint8{1, 2, 3, 4, 5})), []uint8{1, 2, 3, 4, 5}))
	fmt.Println(reflect.DeepEqual(Compact_decode(Compact_encode([]uint8{0, 1, 2, 3, 4, 5})), []uint8{0, 1, 2, 3, 4, 5}))
	fmt.Println(reflect.DeepEqual(Compact_decode(Compact_encode([]uint8{0, 15, 1, 12, 11, 8, 16})), []uint8{0, 15, 1, 12, 11, 8}))
	fmt.Println(reflect.DeepEqual(Compact_decode(Compact_encode([]uint8{15, 1, 12, 11, 8, 16})), []uint8{15, 1, 12, 11, 8}))
}

func (mpt *MerklePatriciaTrie) CreateTestMpt() error {
	mpt.db = make(map[string]Node)

	/*---------- G ------------------ */
	// 0: Null, 1: Branch, 2: Ext or Leaf
	flagValueNodeG := Flag_value{
		encoded_prefix: Compact_encode([]uint8{5, 16}),
		value:          "coin",
	}
	nodeG := Node{
		node_type:  2, //Leaf
		flag_value: flagValueNodeG,
	}
	hashNodeG := nodeG.Hash_node()
	mpt.db[hashNodeG] = nodeG
	fmt.Println("Compact_encode([]uint8{5, 16}:", Compact_encode([]uint8{5, 16}))
	/*---------- F ------------------ */
	flagValueNodeF := Flag_value{
		encoded_prefix: nil,
		value:          "",
	}
	nodeF := Node{
		node_type:    1, //Branch
		flag_value:   flagValueNodeF,
		branch_value: [17]string{"", "", "", "", "", "", hashNodeG, "", "", "", "", "", "", "", "", "", "puppy"},
	}
	hashNodeF := nodeF.Hash_node()
	mpt.db[hashNodeF] = nodeF
	/*---------- E ------------------ */
	flagValueNodeE := Flag_value{
		encoded_prefix: Compact_encode([]uint8{7}),
		value:          hashNodeF,
	}
	nodeE := Node{
		node_type:  2, //Extension
		flag_value: flagValueNodeE,
	}
	hashNodeE := nodeE.Hash_node()
	mpt.db[hashNodeE] = nodeE
	/*---------- J ------------------ */
	flagValueNodeJ := Flag_value{
		encoded_prefix: nil ,
		value:          "book",
	}	
	nodeJ := Node{
		node_type:  2, //Extension
		flag_value: flagValueNodeJ,
	}
	hashNodeJ := nodeJ.Hash_node()
	mpt.db[hashNodeJ] = nodeJ
	/*---------- H ------------------ */
	flagValueNodeH := Flag_value{
		encoded_prefix: nil,
		value:          "",
	}
	nodeH := Node{
		node_type:    1, //Branch
		flag_value:   flagValueNodeH,
		branch_value: [17]string{"", "", hashNodeJ, "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
	}
	hashNodeH := nodeH.Hash_node()
	mpt.db[hashNodeH] = nodeH
	/*---------- D ------------------ */
	flagValueNodeD := Flag_value{
		encoded_prefix: nil,
		value:          "",
	}
	nodeD := Node{
		node_type:    1, //Branch
		flag_value:   flagValueNodeD,
		branch_value: [17]string{"", "", hashNodeH, "", "", "", hashNodeE, "", "", "", "", "", "", "", "", "", "verb"},
	}
	hashNodeD := nodeD.Hash_node()
	mpt.db[hashNodeD] = nodeD

 	/*---------- B ------------------ */
	flagValueNodeB := Flag_value{
		encoded_prefix: Compact_encode([]uint8{6, 15}),
		value:          hashNodeD,
	}
	nodeB := Node{
		node_type:  2, //Extension
		flag_value: flagValueNodeB,
	}
	hashNodeB := nodeB.Hash_node()
	mpt.db[hashNodeB] = nodeB
	/*---------- C ------------------ */
	flagValueNodeC := Flag_value{
		encoded_prefix: Compact_encode([]uint8{6, 15, 7, 2, 7, 3, 6, 5, 16}),
		value:          "stallion",
	}
	nodeC := Node{
		node_type:  2, //Leaf
		flag_value: flagValueNodeC,
	}
	hashNodeC := nodeC.Hash_node()
	mpt.db[hashNodeC] = nodeC
	/*---------- A ------------------ */
	flagValueNodeA := Flag_value{
		encoded_prefix: nil,
		value:          "",
	}
	nodeA := Node{
		node_type:    1, //Branch
		flag_value:   flagValueNodeA,
		branch_value: [17]string{"", "", "", "", hashNodeB, "", "", "", hashNodeC, "", "", "", "", "", "", "", ""},
	}
	hashNodeA := nodeA.Hash_node()
	mpt.db[hashNodeA] = nodeA

	/*---------- Root  ------------------ */
	flagValueRoot := Flag_value{
		encoded_prefix: Compact_encode([]uint8{6}),
		value:          hashNodeA,
	}
	nodeRoot := Node{
		node_type:  2, //Extension Root
		flag_value: flagValueRoot,
	}
	hashNodeRoot := nodeRoot.Hash_node()
	mpt.db[hashNodeRoot] = nodeRoot
	mpt.root = hashNodeRoot
	return errors.New("Problem occured while creating Root Node")
}
