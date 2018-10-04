package core

import (
	"flag"
	"fmt"
		"log"
	"os"
	"strconv"
)

type Client struct {
	BC *BlockChain
}

func (cli *Client) Help() {
	fmt.Println("please send the args:")
	fmt.Println("-add  -data   block_data	add the block to the block chain")
	fmt.Println("-list				print all of the block in the block chain")
}

func (cli *Client) ValidateArgs() {
	if len(os.Args) < 2 {
		cli.Help()
		os.Exit(1)
	}
}

func (cli *Client) Run() {
	cli.ValidateArgs()

	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)

	addData := addCmd.String("data", "", "Block Data")

	switch os.Args[1] {
	case "add":
		err := addCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "list":
		err := listCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.Help()
		os.Exit(1)
	}
	if addCmd.Parsed() {
		if *addData == "" {
			addCmd.Usage()
			os.Exit(1)
		}
		cli.addBlock(*addData)
	}
	if listCmd.Parsed() {
		cli.PrintBlocks()
	}
}
func (cli *Client) addBlock(data string) {
	cli.BC.AddBlock(data)
	fmt.Println("\nSuccess!")
}

func (cli *Client) PrintBlocks() {
	bci := cli.BC.Iterator()
	for {
		block := bci.Next()
		fmt.Printf("Prive hash :%x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := NewProofOfWork(block)
		fmt.Printf("pow:%s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}
