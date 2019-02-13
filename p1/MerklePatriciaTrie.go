package p1

import (
	"encoding/hex"
	"errors"
	"fmt"
	"golang.org/x/crypto/sha3"
	"reflect"
	"strings"
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
const NoChild = uint8(60)
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

/* IsEqualPath
* To compare 2 path and return true if those are the same
*@ input: path1 []uint8, path2 []uint8
*@ output: bool
*/
func IsEqualPath(pathA []uint8, pathB []uint8) bool {
	if pathA == nil {
		return pathB == nil
	}
	if pathB == nil {
		return false
	}
	if len(pathA) != len(pathB) {
		return false
	}
	for i := 0; i < len(pathB); i++ {
		if pathA[i] != pathB[i] {
			return false
		}
	}
	return true
}

/* GetMatchPrefix
* To find match character in path
*@ input: pathA []uint8, pathB []uint8
*@ output: []uint8
*/
func GetMatchPrefix(pathA []uint8, pathB []uint8) []uint8 {
	minLength := len(pathA)
	if minLength > len(pathB) {
		minLength = len(pathB)
	}
	matchPrefix := []uint8{}
	for i := 0; i < minLength; i++ {
		if pathA[i] == pathB[i] {
			matchPrefix = append(matchPrefix, pathA[i])
		} else {
			break
		}
	}
	return matchPrefix
}

func (mpt *MerklePatriciaTrie) UpdateHashValues(nodeChain []Node, child uint8, new_value string) {
	if len(nodeChain) == 0 {
		return
	}

	childrenIndexes := make([]uint8, len(nodeChain))
	for i := 0; i < len(nodeChain)-1; i++ {
		childrenIndexes[i] = NoChild
		if nodeChain[i].IsBranch() {
			nextNodeHash := nodeChain[i+1].hash_node()
			childrenIndexes[i] = nodeChain[i].GetBranchIndex(nextNodeHash)
		}
	}
	childrenIndexes[len(childrenIndexes)-1] = child

	newValue := new_value
	for i := len(nodeChain) - 1; i >= 0; i-- {
		node := nodeChain[i]
		delete(mpt.db, node.hash_node())
		if node.IsLeaf() || node.IsExtension() {
			node.flag_value.value = newValue
		} else {
			node.branch_value[childrenIndexes[i]] = newValue
		}
		newValue = node.hash_node()
		mpt.db[newValue] = node
	}
	mpt.root = newValue
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
	path := HexConverter(key)
	if len(path) == 0 {
		if mpt.root == "" {
			//Empty MPT
			rootNode := NewLeafNodeWithValue(path, new_value)
			mpt.root = rootNode.hash_node()
			mpt.db[mpt.root] = rootNode
		} else {
			rootNode := mpt.GetHashNode(mpt.root)
			if rootNode.IsBranch() {
				//Root is Branch
				delete(mpt.db, mpt.root)
				rootNode.branch_value[16] = new_value
				mpt.root = rootNode.hash_node()
				mpt.db[mpt.root] = rootNode
			} else {
				// Root is Ext or Leaf
				newRootNode := mpt.MergeLeafExt(rootNode, path, new_value)
				mpt.root = newRootNode.hash_node()
			} 
		}
	} else if mpt.root == "" {
		rootNode := NewLeafNodeWithValue(path, new_value)
		mpt.root = rootNode.hash_node()
		mpt.db[mpt.root] = rootNode
	} else {
		nodePath, remainingPath := mpt.GetNodePath(path, mpt.root)
		if len(nodePath) == 0 {
			rootNode := mpt.GetHashNode(mpt.root)
			if rootNode.IsBranch() {
				// Root is Branch node
				childNode := NewLeafNodeWithValue(path[1:], new_value)
				childHash := childNode.hash_node()
				mpt.db[childHash] = childNode
				rootNode.branch_value[path[0]] = childHash
				delete(mpt.db, mpt.root)
				mpt.root = rootNode.hash_node()
				mpt.db[mpt.root] = rootNode
			} else {
				newRootNode := mpt.MergeLeafExt(rootNode, remainingPath, new_value)
				mpt.root = newRootNode.hash_node()
			} 
		} else {
				lastPrefixNode := nodePath[len(nodePath)-1]
				if lastPrefixNode.IsBranch() {
					if len(remainingPath) == 0 {
						mpt.UpdateHashValues(nodePath, 16, new_value)
					} else if lastPrefixNode.branch_value[remainingPath[0]] == "" {
						//Branch Node in Leaf
						newLeafNode := NewLeafNodeWithValue(remainingPath[1:], new_value)
						newLeafNodeHash := newLeafNode.hash_node()
						mpt.db[newLeafNodeHash] = newLeafNode
						mpt.UpdateHashValues(nodePath, remainingPath[0], newLeafNodeHash)
					} else {
						// Find branch and child
						childNode := mpt.GetHashNode(lastPrefixNode.branch_value[remainingPath[0]])
						newChildNode := mpt.MergeLeafExt(childNode, remainingPath[1:], new_value)
						mpt.UpdateHashValues(nodePath, remainingPath[0], newChildNode.hash_node())
					}
				} else {
					lastPrefixNodeHash := lastPrefixNode.hash_node()
					newLastPrefixNode := mpt.MergeLeafExt(lastPrefixNode,append(compact_decode(lastPrefixNode.flag_value.encoded_prefix), remainingPath...),
						new_value)
					if len(nodePath) == 1 {
						mpt.root = newLastPrefixNode.hash_node()
					} else {
						parentNode := nodePath[len(nodePath)-2]
						childIndex := parentNode.GetBranchIndex(lastPrefixNodeHash)
						mpt.UpdateHashValues(nodePath[:len(nodePath)-1], childIndex, newLastPrefixNode.hash_node())
					}
				}
		}
	}
}

/* Delete
* To Delete value from MPT
* Return empty string
*@ input: key string
*@ output: value string, errorMsg error
 */
func (mpt *MerklePatriciaTrie) Delete(key string) (string, error) {
	path := HexConverter(key)
	nodePath, remainingPath := mpt.GetNodePath(path, mpt.root)
	// No key found
	if (len(nodePath) == 0) || (len(remainingPath) > 0) {
		return "", errors.New("path_not_found")
	}

	lastPathNode := &nodePath[len(nodePath)-1]
	if lastPathNode.IsLeaf() {
		// Delete the leaf, update the parent, and then rebalance the trie
		// around the parent, if necessary.
		lastPathNodeHashCode := lastPathNode.hash_node()
		delete(mpt.db, lastPathNodeHashCode)

		if len(nodePath) == 1 {
			// We deleted the root node.
			mpt.root = ""
		} else {
			// Update the parent (it must be a branch node).
			parentNode := nodePath[len(nodePath)-2]
			// CheckState(parentNode.IsBranch(), fmt.Sprintf("FATAL: Find extension with Leaf child"))
			mpt.RebalanceTriePath(nodePath[:len(nodePath)-1],
			parentNode.GetBranchIndex(lastPathNodeHashCode))
		}
		return "", nil
	}

	// CheckState(lastPathNode.IsBranch(),
	// 	fmt.Sprintf("GetNodePath returned a path that does not end in a leaf or a branch: %v",
	// 		nodePath))
	if lastPathNode.branch_value[16] == "" {
		// We found an exact path in the trie for the given key
		// (because of the trie structure), but the key is not in the trie.
		return "", errors.New("path_not_found")
	}

	mpt.RebalanceTriePath(nodePath, 16)
	return "", nil
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
func (mpt *MerklePatriciaTrie) FindLeafNodeToDelete(node Node, searchPath []uint8, s stack) (string, error, []uint8, uint8, Node, stack) {
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

/* GetNodePath
* To find the path of specific node
*@ input: key string, new_value string, s stack
*@ output: None
 */
func (mpt *MerklePatriciaTrie) GetNodePath(path []uint8, nodeHash string) ([]Node, []uint8) {
	node := mpt.GetHashNode(nodeHash)
	if node.IsLeaf() {
		nodePath := compact_decode(node.flag_value.encoded_prefix)
		if IsPathPrefix(nodePath, path) {
			return []Node{node}, path[len(nodePath):]
		}
		return []Node{}, path
	}
	if node.IsExtension() {
		extensionPath := compact_decode(node.flag_value.encoded_prefix)
		if !IsPathPrefix(extensionPath, path) {
			return []Node{}, path
		}
		recNodePath, recRemainingPath := mpt.GetNodePath(path[len(extensionPath):], node.flag_value.value)
		return append([]Node{node}, recNodePath...), recRemainingPath
	}
	// CheckState(node.IsBranch(), fmt.Sprintf("FATAL: Found node of unknown type"))
	if (len(path) == 0) || (node.branch_value[path[0]] == "") {
		return []Node{node}, path
	}
	recNodePath, recRemainingPath := mpt.GetNodePath(path[1:], node.branch_value[path[0]])
	if len(recNodePath) > 0 {
		return append([]Node{node}, recNodePath...), recRemainingPath
	}
	return []Node{node}, path
}

/* GetBranchIndex
* To generate the branch index with value
*@input:childHash string
*@oututp:i uint8
*/
func (node *Node) GetBranchIndex(childHash string) uint8 {

	for i := uint8(0); i < uint8(16); i++ {
		if node.branch_value[i] == childHash {
			return i
		}
	}
	return NoChild
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

/* NewLeafNodeWithValue
* To Print Root
*@ input: path []uint8, value string
*@ output: Node
 */
func NewLeafNodeWithValue(path []uint8, value string) Node {
	return Node{
		node_type:    2,
		branch_value: [17]string{},
		flag_value: Flag_value{
			encoded_prefix: compact_encode(append(path, 16)),
			value:          value,
		},
	}
}

/* GetHashNode
* To pass hash value to DB and get the node
*@ input: hash string
*@ output: n Node
*/
func (mpt *MerklePatriciaTrie) GetHashNode(hash string) Node {
	n, _ := mpt.db[hash]
	return n
}


/* IsBranch
*
* To check if node is Branch node
*@ input: node *Node
*@ output: bool
 */
func (node *Node) IsBranch() bool {
	return node.node_type == 1
}

/* IsLeaf
*
* To check if node is Leaf node
*@ input: node *Node
*@ output: bool
 */
func (node *Node) IsLeaf() bool {
	return (node.node_type == 2) && (node.flag_value.encoded_prefix[0]>>5 == 1)
}

/* IsExtension
* To check for extention node
*@input:node *Node
*@oututp:bool
*/
func (node *Node) IsExtension() bool {
	return (node.node_type == 2) && (node.flag_value.encoded_prefix[0]>>5 == 0)
}

/* hash_node
*
* To convert node as hashNode (HashValue)
*@ input: node *Node
*@ output: errorMessage string
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
func test_compact_encode() {
	fmt.Println("-------------------Test_compact_encode-------------------")
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{1, 2, 3, 4, 5})), []uint8{1, 2, 3, 4, 5}))
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{0, 1, 2, 3, 4, 5})), []uint8{0, 1, 2, 3, 4, 5}))
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{0, 15, 1, 12, 11, 8, 16})), []uint8{0, 15, 1, 12, 11, 8}))
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{15, 1, 12, 11, 8, 16})), []uint8{15, 1, 12, 11, 8}))
}

/* Test_Insert_Get
*
* To Test Insert and Get functions
*/
func Test_Insert_Get_Delete() {
	fmt.Println("-------------------Test_Insert_Get-------------------")
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

	
	fmt.Println("-------------------Test_Get_Delete-------------------")
	mpt.Delete("a")
	mpt.Delete("ab")
	mpt.Delete("acb")
	mpt.Delete("bfge")
	mpt.Delete("c")
	mpt.Delete("r")

	value7, _ := mpt.Get("a")
	value8, _ := mpt.Get("ab")
	value9, _ := mpt.Get("acb")
	value10, _ := mpt.Get("bfge")
	value11, _ := mpt.Get("c")
	value12, _ := mpt.Get("r")

	fmt.Println(reflect.DeepEqual("", value7))
	fmt.Println(reflect.DeepEqual("", value8))
	fmt.Println(reflect.DeepEqual("", value9))
	fmt.Println(reflect.DeepEqual("", value10))
	fmt.Println(reflect.DeepEqual("", value11))
	fmt.Println(reflect.DeepEqual("", value12))
	
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
	fmt.Println("-------------------Test_Get-------------------")
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
	fmt.Println("-------------------Test_Get-------------------")
	fmt.Println(reflect.DeepEqual("pie", value7))
	fmt.Println(reflect.DeepEqual("apple", value8))
}


func (node *Node) String() string {
	str := "empty string"
	switch node.node_type {
	case 0:
		str = "[Null Node]"
	case 1:
		str = "Branch["
		for i, v := range node.branch_value[:16] {
			str += fmt.Sprintf("%d=\"%s\", ", i, v)
		}
		str += fmt.Sprintf("value=%s]", node.branch_value[16])
	case 2:
		encoded_prefix := node.flag_value.encoded_prefix
		node_name := "Leaf"
		if is_ext_node(encoded_prefix) {
			node_name = "Ext"
		}
		ori_prefix := strings.Replace(fmt.Sprint(compact_decode(encoded_prefix)), " ", ", ", -1)
		str = fmt.Sprintf("%s<%v, value=\"%s\">", node_name, ori_prefix, node.flag_value.value)
	}
	return str
}

func node_to_string(node Node) string {
	return node.String()
}

func (mpt *MerklePatriciaTrie) Initial() {
	mpt.root = ""
	mpt.db = make(map[string]Node)
}

func is_ext_node(encoded_arr []uint8) bool {
	return encoded_arr[0] / 16 < 2
}

func TestCompact() {
	test_compact_encode()
}

func (mpt *MerklePatriciaTrie) String() string {
	content := fmt.Sprintf("ROOT=%s\n", mpt.root)
	for hash := range mpt.db {
		content += fmt.Sprintf("%s: %s\n", hash, node_to_string(mpt.db[hash]))
	}
	return content
}

func (mpt *MerklePatriciaTrie) Order_nodes() string {
	raw_content := mpt.String()
	content := strings.Split(raw_content, "\n")
	root_hash := strings.Split(strings.Split(content[0], "HashStart")[1], "HashEnd")[0]
	queue := []string{root_hash}
	i := -1
	rs := ""
	cur_hash := ""
	for len(queue) != 0 {
		last_index := len(queue) - 1
		cur_hash, queue = queue[last_index], queue[:last_index]
		i += 1
		line := ""
		for _, each := range content {
			if strings.HasPrefix(each, "HashStart" + cur_hash + "HashEnd") {
				line = strings.Split(each, "HashEnd: ")[1]
				rs += each + "\n"
				rs = strings.Replace(rs, "HashStart" + cur_hash + "HashEnd", fmt.Sprintf("Hash%v", i),  -1)
			}
		}
		temp2 := strings.Split(line, "HashStart")
		flag := true
		for _, each := range temp2 {
			if flag {
				flag = false
				continue
			}
			queue = append(queue, strings.Split(each, "HashEnd")[0])
		}
	}
	return rs
}








 




