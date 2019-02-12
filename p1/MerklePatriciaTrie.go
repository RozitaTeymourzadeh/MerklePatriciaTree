package p1

import (
	"encoding/hex"
	"errors"
	"fmt"
	"reflect"

	"golang.org/x/crypto/sha3"
)

/*-------------------------STRUCT---------------------------------------------------*/
/* Struct data structure for variables
/*-------------------------STRUCT---------------------------------------------------*/

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

/*-------------------------SERVICE---------------------------------------------------*/
/* Service functions act as helper function to others master functions
/*-------------------------SERVICE---------------------------------------------------*/

/* compact_decode
* To reverse compact encode value
*@ input: encoded_arr []uint8
*@ output: decoded_arr []uint8
 */
func compact_decode(encoded_arr []uint8) []uint8 {

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

/* compact_encode
* To encode hex Key
*@ input: hex_array []uint8
*@ output: result []uint8
 */
func compact_encode(hex_array []uint8) []uint8 {

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

/* GetNodeType
* To decide on Extension and Leaf type
*@ input: node Node
*@ output: nodeType uint8
 */
func GetNodeType(node Node) uint8 {

	var nodeType uint8
	encoded_prefix := node.flag_value.encoded_prefix[0]
	nodeType = encoded_prefix / 16
	return nodeType
}

/* HexConverter
* To convert key into Hex value
*@ input: Key string
*@ output: hex_array []uint8
 */
func HexConverter(key string) []uint8 {
	var hex_array []uint8

	for i := 0; i < len(key); i++ {
		hex_array = append(hex_array, key[i]/16)
		hex_array = append(hex_array, key[i]%16)
	}
	return hex_array
}

/* InitializeMpt
* To initialize mpt
*@ input: None
*@ output: pointerto tree
 */
func InitializeMpt() *MerklePatriciaTrie {
	db := make(map[string]Node)
	root := ""
	return &MerklePatriciaTrie{db, root}
}

/* InitializeStack
* To initialize stack
*@ input: None
*@ output: s stack
 */
func InitializeStack() stack {
	s := make(stack, 0)
	return s
}

/* Push
* To push into stack
*@ input: v None
*@ output: s stack
 */
func (s stack) Push(v Node) stack {
	return append(s, v)
}

/* Pop
* To Pop from stack
*@ input: n Node, s stack
*@ output: n Node,s stack
 */
func (s stack) Pop() (stack, Node) {
	l := len(s)
	if l == 0 {
		fmt.Println("Stack is Empty!!")
		return s, Node{}
	}
	return s[:l-1], s[l-1]
}

/* EqualArray
* To compare 2 arrays and return no. of match index
*@ input: a []uint8, b []uint8`
*@ output: j int , remainPath []uint8
 */
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

/*-------------------------MASTER---------------------------------------------------*/
/* Master Functions
/*-------------------------MASTER---------------------------------------------------*/

/* Get
* To receive key and return associated value with key
* It return empty string if not find
*@ input: key string
*@ output: value string, errorMsg error
 */
func (mpt *MerklePatriciaTrie) Get(key string) (string, error) {
	var DEBUG int
	DEBUG = 0
	var value string
	var errorMsg error
	var nodeType uint8
	nodeType = 0
	if key == "" {
		fmt.Println("FATAl: Invalid KEY ...")
		return "", errors.New("FATAl: Invalid KEY ...")
	}
	currentNode := mpt.db[mpt.root]
	if currentNode.node_type == 2 {
		nodeType = GetNodeType(currentNode)
	}
	searchPath := HexConverter(key)
	for len(searchPath) != 0 && value == "" && errorMsg == nil {
		// if node_type is NULL, Branch, Leaf/Ext
		if nodeType < 2 {
			//Extention or Branch
			valueFln, errorMsgFln, remainPath, nodeTypeFln, nextNode := mpt.FindLeafNode(currentNode, searchPath)
			searchPath = remainPath
			currentNode = nextNode
			errorMsg = errorMsgFln
			nodeType = nodeTypeFln
			value = valueFln
		} else {
			//Leaf
			valueFlv, errorMsgFlv, remainPath := FindLeafValue(currentNode, searchPath)
			if DEBUG == 1 {
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

/* Insert
* To Insert value into MPT
* 
*@ input: key string, new_value string
*@ output: none
 */
func (mpt *MerklePatriciaTrie) Insert(key string, new_value string) {
	// TODO
}

/* Delete
* To Delete value from MPT
* Return empty string
*@ input: key string
*@ output: value, errorMsg
 */
func (mpt *MerklePatriciaTrie) Delete(key string) (string, error) {
	var value string

	if key == "" {
		fmt.Println("FATAl: Invalid KEY ...")
		return "", errors.New("path_not_found")
	}
	value, err, s := mpt.GetToDelete(key)
	fmt.Println("Slack length is:", len(s))
	s, p := s.Pop()
	fmt.Println(p)
	s, p = s.Pop()
	fmt.Println(p)
	s, p = s.Pop()
	fmt.Println(p)
	s, p = s.Pop()
	fmt.Println(p)

	fmt.Println(err)
	return value, errors.New("path_not_found")
}


/*-------------------------GET HELPER---------------------------------------------------*/
/* SubFunction of Get Master Function
/*-------------------------GET HELPER---------------------------------------------------*/

/* FindLeafValue
* To find value in Leaf node
* The pointer is in leaf node when we are here
*@ input: node Node, searchPath []uint8
*@ output: value string, errorMsg error, remainPath []uint8
 */
 func FindLeafValue(node Node, searchPath []uint8) (string, error, []uint8) {
	currentPath := compact_decode(node.flag_value.encoded_prefix)
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

/* FindLeafNode
* To find Leaf Node 
*@ input: node Node, searchPath []uint8
*@ output: string, error, []uint8, uint8, Node
 */
func (mpt *MerklePatriciaTrie) FindLeafNode(node Node, searchPath []uint8) (string, error, []uint8, uint8, Node) {
	var value string
	var remainPath []uint8
	var nodeType uint8
	var nextNode Node
	var currentPath []uint8
	var index uint8

	if node.node_type == 2 {
		// Extension or Leaf
		currentPath = compact_decode(node.flag_value.encoded_prefix)
	} else if node.node_type == 1 {
		// Branch
		branchCurrentPath := node.branch_value
		for index = 0; index < uint8(len(branchCurrentPath)); index++ {
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
	if node.branch_value[searchPath[0]] != "" && node.node_type == 1 {
		nextNode = mpt.db[node.branch_value[searchPath[0]]]
		if nextNode.node_type == 1 && len(searchPath) == 1 {
			value = nextNode.branch_value[16]
		} else if nextNode.node_type == 2 {
			value = nextNode.flag_value.value
		}

		return value, nil, remainPath, nodeType, nextNode
	}
	// if whole match path and remaining value
	if matchedIndex+1 == len(currentPath) && len(remainPath) != 0 {
		nextBranchNode := mpt.db[node.flag_value.value]
		if nextBranchNode.branch_value[remainPath[0]] != "" {
			nextNode = mpt.db[nextBranchNode.branch_value[remainPath[0]]]
			if nextNode.node_type == 1 {
				// Node is Branch
				remainPath = remainPath[1:]
			} else if nextNode.node_type == 2 {
				// Leaf or Extension
				nodeType = GetNodeType(nextNode)
				if nodeType < 2 {
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
	return "", nil, remainPath, nodeType, nextNode
}


/*-------------------------DELETE HELPER---------------------------------------------------*/
/* SubFunction of Get Master Function
/*-------------------------DELETE HELPER---------------------------------------------------*/

/* GetToDelete 
* Type: (Optional)
* To find key value pair for delete purpose
* It return empty string if not find
*@ input: key string
*@ output: value string, errorMsg error, s stack
 */
func (mpt *MerklePatriciaTrie) GetToDelete(key string) (string, error, stack) {
	var value string
	var errorMsg error
	var nodeType uint8
	nodeType = 0
	s := InitializeStack()
	if key == "" {
		fmt.Println("FATAl: Invalid KEY ...")
		return "", errors.New("path_not_found"), nil
	}
	currentNode := mpt.db[mpt.root]
	if currentNode.node_type != 0 {
		s = s.Push(currentNode)
	}

	if currentNode.node_type == 2 {
		nodeType = GetNodeType(currentNode)
	}
	searchPath := HexConverter(key)
	for len(searchPath) != 0 && value == "" && errorMsg == nil {
		// if node_type is NULL, Branch, Leaf/Ext
		if nodeType < 2 {
			//Extention or Branch
			valueFln, errorMsgFln, remainPath, nodeTypeFln, nextNodeFln, sFln := mpt.FindLeafNodeToDelete(currentNode, searchPath, s)

			searchPath = remainPath
			currentNode = nextNodeFln
			errorMsg = errorMsgFln
			nodeType = nodeTypeFln
			value = valueFln
			s = sFln
		} else {
			//Leaf
			valueFlv, errorMsgFlv, _ := FindLeafValue(currentNode, searchPath)
			value = valueFlv
			errorMsg = errorMsgFlv
		}
	}
	return value, errorMsg, s
}

/* FindLeafNodeToDelete
* Type: (Optional)
* To find Leaf Node for Delete purpose
*@ input: node Node, searchPath []uint8
*@ output: string, error, []uint8, uint8, Node
 */
func (mpt *MerklePatriciaTrie) FindLeafNodeToDelete(node Node, searchPath []uint8, s stack) (string, error, []uint8, 		uint8, Node, stack) {
	var value string
	var remainPath []uint8
	var nodeType uint8
	var nextNode Node
	var currentPath []uint8
	var index uint8

	if node.node_type == 2 {
		// Extension or Leaf
		currentPath = compact_decode(node.flag_value.encoded_prefix)
	} else if node.node_type == 1 {
		// Branch
		branchCurrentPath := node.branch_value
		for index = 0; index < uint8(len(branchCurrentPath)); index++ {
			if branchCurrentPath[index] != "" {
				currentPath = append(currentPath, index)
			}
		}
	} else {
		// Null
		fmt.Println("Null Branch!!")
		return "", errors.New("path_not_found"), remainPath, nodeType, nextNode, s
	}
	matchedIndex, remainPath := EqualArray(currentPath, searchPath)
	// if whole match, return Branch Node Value if current Node is Ext/Leaf
	if matchedIndex+1 == len(currentPath) && len(remainPath) == 0 && node.node_type == 2 {
		nextBranchNode := mpt.db[node.flag_value.value]
		s = s.Push(nextBranchNode)
		value = nextBranchNode.branch_value[16]
		return value, nil, remainPath, nodeType, nextBranchNode, s
	}
	// if whole match, if current Node is Branch
	if node.branch_value[searchPath[0]] != "" && node.node_type == 1 {
		nextNode = mpt.db[node.branch_value[searchPath[0]]]
		s = s.Push(nextNode)
		if nextNode.node_type == 2 {
			value = nextNode.flag_value.value
		}
		return value, nil, remainPath, nodeType, nextNode, s
	}
	// if whole match path and remaining value
	if matchedIndex+1 == len(currentPath) && len(remainPath) != 0 {
		nextBranchNode := mpt.db[node.flag_value.value]
		s = s.Push(nextBranchNode)
		if nextBranchNode.branch_value[remainPath[0]] != "" {
			nextNode = mpt.db[nextBranchNode.branch_value[remainPath[0]]]
			s = s.Push(nextNode)
			if nextNode.node_type == 1 {
				// Node is Branch
				remainPath = remainPath[1:]
			} else if nextNode.node_type == 2 {
				// Leaf or Extension
				nodeType = GetNodeType(nextNode)
				if nodeType < 2 {
					//Extension
					remainPath = remainPath[1:]
				} else {
					//Leaf
					remainPath = remainPath[1:]
				}
			}
		} else {
			// if no match path
			return "", errors.New("path_not_found"), remainPath, nodeType, nextNode, s
		}
	}
	return "", nil, remainPath, nodeType, nextNode, s
}

/*-------------------------INSERT HELPER---------------------------------------------------*/
/* SubFunction of Insert Master Function
/*-------------------------INSERT HELPER---------------------------------------------------*/

/* InsertRoot
* Type: (Optional)
* To Insert value on root
*@ input: key string, new_value string, s stack
*@ output: None
 */
func (mpt *MerklePatriciaTrie) InsertRoot(key string, new_value string, s stack) {
	// first Insert
	if len(s) == 0 {
		hex_array := HexConverter(key)
		hex_array = append(hex_array, 16)
		flagValue := Flag_value{
			encoded_prefix: compact_encode(hex_array),
			value:          new_value,
		}
		newNode := Node{
			node_type:  2,
			flag_value: flagValue,
		}
		hashedNode := newNode.hash_node()
		mpt.db[hashedNode] = newNode
		mpt.root = hashedNode
		return
	}
}


/*-------------------------NODE HELPER---------------------------------------------------*/
/* Node accessories
/*-------------------------NODE HELPER---------------------------------------------------*/

/* GetRootNode
* Type: (Optional)
* To Print Root
*@ input: None
*@ output: None
 */

func (mpt *MerklePatriciaTrie) GetRootNode() {
	root := mpt.db[mpt.root]
	fmt.Println("Nodetype: ", root.node_type)
	fmt.Println("Branch:", root.branch_value)
	fmt.Println("Prefix:", root.flag_value.encoded_prefix)
	fmt.Println("Value: ", root.flag_value.value)
}

/* hash_node
* 
* To convert node as hashNode (HashValue)
*@ input: node *Node
*@ output: str string
 */
 func (node *Node) hash_node() string {
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


/*-------------------------TEST---------------------------------------------------*/
/* To Test the project (Functions and Subfunction)
/*-------------------------TEST---------------------------------------------------*/

/* CreateTestMpt
* 
* To create dummy Mpt
*/
func (mpt *MerklePatriciaTrie) CreateTestMpt() error {
	mpt.db = make(map[string]Node)

	/*---------- G ------------------ */
	// 0: Null, 1: Branch, 2: Ext or Leaf
	flagValueNodeG := Flag_value{
		encoded_prefix: compact_encode([]uint8{5, 16}),
		value:          "coin",
	}
	nodeG := Node{
		node_type:  2, //Leaf
		flag_value: flagValueNodeG,
	}
	hashNodeG := nodeG.hash_node()
	mpt.db[hashNodeG] = nodeG

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
	hashNodeF := nodeF.hash_node()
	mpt.db[hashNodeF] = nodeF
	/*---------- E ------------------ */
	flagValueNodeE := Flag_value{
		encoded_prefix: compact_encode([]uint8{7}),
		value:          hashNodeF,
	}
	nodeE := Node{
		node_type:  2, //Extension
		flag_value: flagValueNodeE,
	}
	hashNodeE := nodeE.hash_node()
	mpt.db[hashNodeE] = nodeE
	/*---------- J ------------------ */
	flagValueNodeJ := Flag_value{
		encoded_prefix: nil,
		value:          "book",
	}
	nodeJ := Node{
		node_type:  2, //Extension
		flag_value: flagValueNodeJ,
	}
	hashNodeJ := nodeJ.hash_node()
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
	hashNodeH := nodeH.hash_node()
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
	hashNodeD := nodeD.hash_node()
	mpt.db[hashNodeD] = nodeD

	/*---------- B ------------------ */
	flagValueNodeB := Flag_value{
		encoded_prefix: compact_encode([]uint8{6, 15}),
		value:          hashNodeD,
	}
	nodeB := Node{
		node_type:  2, //Extension
		flag_value: flagValueNodeB,
	}
	hashNodeB := nodeB.hash_node()
	mpt.db[hashNodeB] = nodeB
	/*---------- C ------------------ */
	flagValueNodeC := Flag_value{
		encoded_prefix: compact_encode([]uint8{6, 15, 7, 2, 7, 3, 6, 5, 16}),
		value:          "stallion",
	}
	nodeC := Node{
		node_type:  2, //Leaf
		flag_value: flagValueNodeC,
	}
	hashNodeC := nodeC.hash_node()
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
	hashNodeA := nodeA.hash_node()
	mpt.db[hashNodeA] = nodeA

	/*---------- Root  ------------------ */
	flagValueRoot := Flag_value{
		encoded_prefix: compact_encode([]uint8{6}),
		value:          hashNodeA,
	}
	nodeRoot := Node{
		node_type:  2, //Extension Root
		flag_value: flagValueRoot,
	}
	hashNodeRoot := nodeRoot.hash_node()
	mpt.db[hashNodeRoot] = nodeRoot
	mpt.root = hashNodeRoot
	return errors.New("Problem occured while creating Root Node")
}

/* CreateTestMpt3
* 
* To create dummy Mpt
*/
func (mpt *MerklePatriciaTrie) CreateTestMpt3() error {
	mpt.db = make(map[string]Node)

	flagValueNodeC := Flag_value{
		encoded_prefix: compact_encode([]uint8{2}),
		value:          "pie",
	}
	nodeC := Node{
		node_type:  2, //Leaf
		flag_value: flagValueNodeC,
	}
	hashNodeC := nodeC.hash_node()
	mpt.db[hashNodeC] = nodeC

	flagValueNodeB := Flag_value{
		encoded_prefix: compact_encode([]uint8{1}),
		value:          "apple",
	}
	nodeB := Node{
		node_type:  2, //Leaf
		flag_value: flagValueNodeB,
	}
	hashNodeB := nodeB.hash_node()
	mpt.db[hashNodeB] = nodeB

	flagValueNodeRootA := Flag_value{
		encoded_prefix: nil,
		value:          "",
	}
	nodeRootA := Node{
		node_type:    1, //Branch
		flag_value:   flagValueNodeRootA,
		branch_value: [17]string{"", "", "", "", "", "", hashNodeB, hashNodeC, "", "", "", "", "", "", "", "", ""},
	}

	hashNodeRootA := nodeRootA.hash_node()
	mpt.db[hashNodeRootA] = nodeRootA
	mpt.root = hashNodeRootA
	return errors.New("Problem occured while creating Root Node")
}

/* Test_compact_encode
* 
* To Test compact_encode()
*/
func Test_compact_encode() {
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{1, 2, 3, 4, 5})), []uint8{1, 2, 3, 4, 5}))
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{0, 1, 2, 3, 4, 5})), []uint8{0, 1, 2, 3, 4, 5}))
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{0, 15, 1, 12, 11, 8, 16})), []uint8{0, 15, 1, 12, 11, 8}))
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{15, 1, 12, 11, 8, 16})), []uint8{15, 1, 12, 11, 8}))
}

/* Test_Insert_Get
* 
* To Test Insert and Get functions
*/
func Test_Insert_Get() {
	mpt := InitializeMpt()
	mpt.Insert("a", "apple")
	mpt.Insert("ab", "banana")
	mpt.Insert("acb", "horse")
	mpt.Insert("bfge", "Dog")
	mpt.Insert("c", "Doggy")
	mpt.Insert("r", "Lucy")
	value1, _ := mpt.Get("a")
	value2, _ := mpt.Get("ab")
	value3, _ := mpt.Get("acb")
	value4, _ := mpt.Get("bfge")
	value5, _ := mpt.Get("c")
	value6, _ := mpt.Get("r")

	fmt.Println(reflect.DeepEqual("apple", value1))
	fmt.Println(reflect.DeepEqual("banana", value2))
	fmt.Println(reflect.DeepEqual("horse", value3))
	fmt.Println(reflect.DeepEqual("Dog", value4))
	fmt.Println(reflect.DeepEqual("Doggy", value5))
	fmt.Println(reflect.DeepEqual("Lucy", value6))
}

/* Test_Get
* 
* To Test Get functions
*/
func (mpt *MerklePatriciaTrie) Test_Get() {
	mpt = InitializeMpt()
	mpt.CreateTestMpt()

	value1, _ := mpt.Get("do")
	value2, _ := mpt.Get("dog")
	value3, _ := mpt.Get("doge")
	value4, _ := mpt.Get("horse")
	value5, _ := mpt.Get("do\"")
	value6, _ := mpt.Get("")

	fmt.Println(reflect.DeepEqual("verb", value1))
	fmt.Println(reflect.DeepEqual("puppy", value2))
	fmt.Println(reflect.DeepEqual("coin", value3))
	fmt.Println(reflect.DeepEqual("stallion", value4))
	fmt.Println(reflect.DeepEqual("book", value5))
	fmt.Println(reflect.DeepEqual("", value6))

	mpt = InitializeMpt()
	mpt.CreateTestMpt3()

	value7, _ := mpt.Get("r")
	value8, _ := mpt.Get("a")

	fmt.Println(reflect.DeepEqual("pie", value7))
	fmt.Println(reflect.DeepEqual("apple", value8))
}
