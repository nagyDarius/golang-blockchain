package cli

import (
	"flag"
	"fmt"
	"mrnagy.com/owlchain/blockchain"
	"os"
	"runtime"
)

type Cli struct {
	Chain *blockchain.Chain
}

func (cli *Cli) PrintUsage() {
	fmt.Println("Usage:")
	fmt.Println(" add -block BLOCK_DATA -- adds a block to the chain")
	fmt.Println(" print -- prints the blockchain")
}

func (cli *Cli) ValidateArgs() {
	if len(os.Args) < 2 {
		cli.printUsageAndExit()
	}
}

func (cli *Cli) AddBlock(data string) {
	cli.Chain.AddBlock(data)
	fmt.Println("Added Block!")
}

func (cli *Cli) PrintChain() {
	iter := cli.Chain.Iterator()

	for {
		block := iter.Next()

		fmt.Printf("\nHash:     %x\n", block.Hash)
		fmt.Printf("Transactions:     %s\n", block.Transactions)
		fmt.Printf("PrevHash: %x\n", block.PrevHash)
		fmt.Printf("Nonce: 	 %d  Valid: %t\n\n", block.Nonce, blockchain.NewProof(block).Validate())

		if len(block.PrevHash) == 0 {
			break
		}
	}
}

func (cli *Cli) Run() {
	cli.ValidateArgs()

	addBlockCmd := flag.NewFlagSet("add", flag.ExitOnError)
	printBlockCmd := flag.NewFlagSet("print", flag.ExitOnError)
	addBlockData := addBlockCmd.String("block", "", "Block Transactions")

	switch os.Args[1] {
	case "add":
		err := addBlockCmd.Parse(os.Args[2:])
		blockchain.Handle(err)
	case "print":
		err := printBlockCmd.Parse(os.Args[2:])
		blockchain.Handle(err)
	default:
		cli.printUsageAndExit()
	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			cli.printUsageAndExit()
		}
		cli.AddBlock(*addBlockData)
	}

	if printBlockCmd.Parsed() {
		cli.PrintChain()
	}
}

func (cli *Cli) printUsageAndExit() {
	cli.PrintUsage()
	runtime.Goexit()
}
