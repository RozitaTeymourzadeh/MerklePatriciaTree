package main

import (
	"MerklePatriciaTree/p1"
	"fmt"
)

func main() {
	fmt.Println("hello world!")
	/* Initialize MPT and stack */
	mpt := p1.InitializeMpt()
	//s := p1.InitializeStack()

	/* Test Create Node */
	//node1 := p1.Node{}
	//node2 := p1.Node{}
	//node3 := p1.Node{}

	//mpt.InsertRoot("a", "apple", s)

	/* Test GetRoot */
	//mpt.GetRootNode()

	/* Test Get */
	mpt.CreateTestMpt()
	// value, errorMsg := mpt.Get("do")
	// fmt.Println("Get value is : ", value, "Get ErrorMsg is :", errorMsg)
	// value, errorMsg = mpt.Get("dog")
	// fmt.Println("Get value is : ", value, "Get ErrorMsg is :", errorMsg)
	// value, errorMsg = mpt.Get("doge")
	// fmt.Println("Get value is : ", value, "Get ErrorMsg is :", errorMsg)
	// value, errorMsg = mpt.Get("horse")
	// fmt.Println("Get value is : ", value, "Get ErrorMsg is :", errorMsg)
	// value, errorMsg = mpt.Get("do\"")
	// fmt.Println("Get value is : ", value, "Get ErrorMsg is :", errorMsg)
	// value, errorMsg = mpt.Get("")
	// fmt.Println("Get value is : ", value, "Get ErrorMsg is :", errorMsg)

	// mpt.CreateTestMpt3()
	// value, errorMsg = mpt.Get("r")
	// fmt.Println("Get value is : ", value, "Get ErrorMsg is :", errorMsg)
	// value, errorMsg = mpt.Get("a")
	// fmt.Println("Get value is : ", value, "Get ErrorMsg is :", errorMsg)

	value, err := mpt.Delete("do")
	fmt.Println(value, err)
	 


	/* Test ArrayEqual */
	//j, remainPath := p1.EqualArray([]uint8{1, 6, 1}, []uint8{1, 6, 1, 1})
	//fmt.Println("Number of simmilar index:", j, remainPath)

	/* Test Stack */
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

	/* Test HexConverter Encoder and Decoder */
	//key := "a"
	//fmt.Println("key is : ", key)
	//fmt.Println("Converted Hex and input to Encoded value is: ", p1.HexConverter(key))
	//fmt.Println("Encoded value of key is : ", p1.Compact_encode(p1.HexConverter(key)))
	//fmt.Println("decoded value of key is : ", p1.Compact_decode(p1.Compact_encode(p1.HexConverter(key))))

	// fmt.Println("Hex : ", p1.HexConverter("do"))
	// fmt.Println("encode : ", p1.Compact_encode(p1.HexConverter("do")))
	// fmt.Println("decoded : ", p1.Compact_decode(p1.Compact_encode(p1.HexConverter("do"))))
	// fmt.Println("encode : ", p1.Compact_encode([]uint8{2, 16}))

	//fmt.Println("Value:", p1.Compact_decode(p1.Compact_encode([]uint8{1, 6, 1})))
	//fmt.Println("Value:", p1.Compact_decode(p1.Compact_encode([]uint8{2, 0, 6, 1})))
	//fmt.Println("Value:", p1.Compact_decode(p1.Compact_encode([]uint8{0, 15, 1, 12, 11, 8, 16})))
	//fmt.Println("Value:", p1.Compact_decode(p1.Compact_encode([]uint8{15, 1, 12, 11, 8, 16})))
	//fmt.Println("Value:", p1.Compact_decode(p1.Compact_encode([]uint8{})))

	//p1.Test_compact_encode()
}
