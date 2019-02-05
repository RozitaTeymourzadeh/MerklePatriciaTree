package main

import (
	"MerklePatriciaTree/p1"
	"fmt"
)

func main() {
	fmt.Println("hello world!")
	fmt.Println("Encoded Value:", p1.Compact_encode([]uint8{1, 2, 3, 4, 5}))
	fmt.Println("Encoded Value:", p1.Compact_encode([]uint8{0, 1, 2, 3, 4, 5}))
	fmt.Println("Encoded Value:", p1.Compact_encode([]uint8{0, 15, 1, 12, 11, 8, 16}))
	fmt.Println("Encoded Value:", p1.Compact_encode([]uint8{15, 1, 12, 11, 8, 16}))

	//%p1.Test_compact_encode()
}
