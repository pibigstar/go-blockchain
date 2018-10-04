package main

import (
	"go-blockchain/src/core"
)

func main() {
	bc := core.NewBlockChain()
	defer bc.Close()
	cli := &core.Client{bc}
	cli.Run()
}
