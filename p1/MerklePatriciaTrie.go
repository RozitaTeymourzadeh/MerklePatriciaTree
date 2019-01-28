package p1

import (
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/sha3"
	"reflect"
)

type Flag_value struct {
	encoded_prefix []uint8
	value string
}

type Node struct {
	node_type int // 0: Null, 1: Branch, 2: Ext or Leaf
	branch_value [17]string
	flag_value Flag_value
}

type MerklePatriciaTrie struct {
	db map[string]Node
	root string
}

func (mpt *MerklePatriciaTrie) Get(key string) string {
	// TODO
	return ""
}

func (mpt *MerklePatriciaTrie) Insert(key string, new_value string) {
	// TODO
}

func (mpt *MerklePatriciaTrie) Delete(key string) {
	// TODO
}

func compact_encode(hex_array []uint8) []uint8 {
	// TODO
	return []uint8{}
}

// If Leaf, ignore 16 at the end
func compact_decode(encoded_arr []uint8) []uint8 {
	// TODO
	return []uint8{}
}

func test_compact_encode() {
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{1, 2, 3, 4, 5})), []uint8{1, 2, 3, 4, 5}))
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{0, 1, 2, 3, 4, 5})), []uint8{0, 1, 2, 3, 4, 5}))
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{0, 15, 1, 12, 11, 8, 16})), []uint8{0, 15, 1, 12, 11, 8}))
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{15, 1, 12, 11, 8, 16})), []uint8{15, 1, 12, 11, 8}))
}

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