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

func GetPrefix(encoded_prefix uint8) uint8 {

	var decoded_prefix uint8
	if encoded_prefix == 0 {
		fmt.Println("Invalid input data for Decoded_prefix")
		return decoded_prefix
	}
	decoded_prefix = encoded_prefix / 16
	return decoded_prefix
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
func EqualArray(a, b []uint8) int {
	var j int
	// var i int
	// i = 0
	j = -1
	if len(a) == 0 || len(b) == 0 {
		return j
	}
	for i, v := range a {
		if v == b[i] {
			j++
		} else {
			break
		}
	}
	fmt.Println("Number of simmilar index:", j)
	return j
}

func (mpt *MerklePatriciaTrie) Get(key string) (string, error) {
	// TODO
	var value string
	hexPath := HexConverter(key)
	currentNode := mpt.db[mpt.root]
	currentPath := Compact_decode(currentNode.flag_value.encoded_prefix)

	firstIndex := currentNode.flag_value.encoded_prefix[0]
	NodeType := GetPrefix(firstIndex)

	matchedIndex := EqualArray(currentPath, hexPath)
	switch NodeType {
	// Extension
	case 0, 1:
		// if whole match, return Branch Node Value
		if matchedIndex == len(currentPath) {
			nextBranchNode := mpt.db[currentNode.flag_value.value]
			value = nextBranchNode.flag_value.value
		} else if matchedIndex == 0 {
			// if 1 match path
			nextBranchNode := mpt.db[currentNode.flag_value.value]
			if nextBranchNode.branch_value[hexPath[1]] != "" {
				nextNode := mpt.db[nextBranchNode.branch_value[hexPath[1]]]
				// if next node is Leaf {}

				//if next node is Extension {}

			}

		}

	// Leaf
	case 2, 3:
		// if whole match, return Leaf Node Value
		if matchedIndex == len(currentPath) {
			value = currentNode.flag_value.value
			return value, nil
		} else {
			// if NOT match, return Leaf Node Value
			return "", errors.New("path_not_found")
		}
	default:
		return "", errors.New("path_not_found")
	}

	return value, errors.New("path_not_found")
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
