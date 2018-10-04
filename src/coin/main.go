package main

import (
	"fmt"
	"go-blockchain/src/core"
)

func main() {

	bc := core.NewBlockChain()
	bc.AddBlock("Send 1 bit to pi")
	bc.AddBlock("Send 2 bit to pi")

	for _, block := range bc.Blocks {
		fmt.Printf("Hash:%x\n", block.Hash)
		fmt.Printf("Data:%s\n", block.Data)
		fmt.Printf("PrevHash:%x\n", block.PrevBlockHash)
	}

}
