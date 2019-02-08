package main

import (
	"MerklePatriciaTree/p1"
	"fmt"
)

func main() {
	fmt.Println("hello world!")
	mpt := p1.InitializeMpt()
	s := p1.InitializeStack()

	//node1 := p1.Node{}
	//node2 := p1.Node{}
	//node3 := p1.Node{}

	mpt.InsertRoot("a", "apple", s)
	mpt.GetRootNode()
	mpt.Get("a")
	// j := p1.EqualArray([]uint8{},[]uint8{1, 6, 1})
	//fmt.Println("Number of simmilar index:",j)
	// s = s.Push(node1)
	// s = s.Push(node2)
	// s = s.Push(node3)

	// s, p := s.Pop()
	// fmt.Println(p)
	// s, p = s.Pop()
	// fmt.Println(p)
	// s, p = s.Pop()
	// fmt.Println(p)
	// s, p = s.Pop()
	// fmt.Println(p)

	//key := "a"
	//fmt.Println("key is : ", key)
	//fmt.Println("Converted Hex and input to Encoded value is: ", p1.HexConverter(key))
	//fmt.Println("Encoded value of key is : ", p1.Compact_encode(p1.HexConverter(key)))
	//fmt.Println("decoded value of key is : ", p1.Compact_decode(p1.Compact_encode(p1.HexConverter(key))))

	//fmt.Println("Hex : ", p1.HexConverter("do"))
	//fmt.Println("encode : ", p1.Compact_encode(p1.HexConverter("do")))
	//fmt.Println("decoded : ", p1.Compact_decode(p1.Compact_encode(p1.HexConverter("do"))))
	//fmt.Println("encode : ", p1.Compact_encode([]uint8{2, 16}))

	//fmt.Println("Value:", p1.Compact_decode(p1.Compact_encode([]uint8{1, 6, 1})))
	//fmt.Println("Value:", p1.Compact_decode(p1.Compact_encode([]uint8{2, 0, 6, 1})))
	//fmt.Println("Value:", p1.Compact_decode(p1.Compact_encode([]uint8{0, 15, 1, 12, 11, 8, 16})))
	//fmt.Println("Value:", p1.Compact_decode(p1.Compact_encode([]uint8{15, 1, 12, 11, 8, 16})))
	//fmt.Println("Value:", p1.Compact_decode(p1.Compact_encode([]uint8{})))

	//p1.Test_compact_encode()

}
