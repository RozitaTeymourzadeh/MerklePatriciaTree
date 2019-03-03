package main

import (
	"fmt"
	//"MerklePatriciaTree/p1"
	"MerklePatriciaTree/p2"
)

func main() {
	fmt.Println("hello world!")

	/* Test Cases  for project 1*/
	// p1.TestCompact()
	// p1.Test1_2()
	// p1.Test3_4()
	// p1.Test5()
	// p1.Test6()



mpt := p2.MerklePatriciaTrie{}
mpt.Initial()
block := p2.Block{}
blockChain := p2.BlockChain{}
blockChain.Initial()
fmt.Printf("mpt:", mpt)
fmt.Printf("block:", block)
fmt.Printf("blockChain:", blockChain)
}
