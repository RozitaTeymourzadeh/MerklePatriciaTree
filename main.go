package main

import (
	"MerklePatriciaTree/p1"
	"fmt"
)

func main() {
	fmt.Println("hello world!")
	p1.InitializeMpt();
	//key := "a"
	//fmt.Println("key is : ", key)
	//fmt.Println("Converted Hex and input to Encoded value is: ", p1.HexConverter(key))
	//fmt.Println("Encoded value of key is : ", p1.Compact_encode(p1.HexConverter(key)))
	//fmt.Println("decoded value of key is : ", p1.Compact_decode(p1.Compact_encode(p1.HexConverter(key))))

	 //fmt.Println("Value:", p1.Compact_decode(p1.Compact_encode([]uint8{1, 6, 1})))
	 //fmt.Println("Value:", p1.Compact_decode(p1.Compact_encode([]uint8{2, 0, 6, 1})))
	 //fmt.Println("Value:", p1.Compact_decode(p1.Compact_encode([]uint8{0, 15, 1, 12, 11, 8, 16})))
	 //fmt.Println("Value:", p1.Compact_decode(p1.Compact_encode([]uint8{15, 1, 12, 11, 8, 16})))
	 //fmt.Println("Value:", p1.Compact_decode(p1.Compact_encode([]uint8{})))

	//p1.Test_compact_encode()

	
}
