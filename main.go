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

	/* Test GetRoot */
	//mpt.GetRootNode()

 /* Test Cases */
		p1.Test_compact_encode()
		mpt.Test_Get()
	 //p1.Test_Insert_Get()
	 
	
 




	
}
