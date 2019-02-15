package main

import (
	"MerklePatriciaTree/p1"
	"fmt"
)

func main() {
	fmt.Println("hello world!")
	/* Initialize MPT */
		mpt := p1.InitializeMpt()

 /* Test Cases */
		  p1.TestCompact()
		 	mpt.Test_Get()
			p1.Test_Insert_Get_Delete()
			p1.TestInsertDeleteGet3()
			p1.TestInsertGet2()
}
