package cli

import (
	"flag"
	"fmt"
	"mrnagy.com/owlchain/blockchain"
	"os"
	"runtime"
)

type Cli struct {
}

func (cli *Cli) PrintUsage() {
	fmt.Println("Usage:")
	fmt.Println(" balance -address ADDRESS -- gets the balance for the address")
	fmt.Println(" create -address ADDRESS -- creates a blockchain for the address")
	fmt.Println(" print -- prints the blockchain")
	fmt.Println(" send -from FROM -to TO -amount AMOUNT -- sends tokens")
}

func (cli *Cli) ValidateArgs() {
	if len(os.Args) < 2 {
		cli.printUsageAndExit()
	}
}

func (cli *Cli) PrintChain() {
	chain := blockchain.ContinueBlockChain("")
	defer chain.Database.Close()
	iter := chain.Iterator()

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

func (cli *Cli) CreateBlockChain(address string) {
	chain := blockchain.InitBlockChain(address)
	chain.Database.Close()
	fmt.Println("Finished!")
}

func (cli *Cli) Balance(address string) {
	chain := blockchain.ContinueBlockChain(address)
	defer chain.Database.Close()

	balance := 0
	unspentOutputs := chain.FindUTXO(address)

	for _, output := range unspentOutputs {
		balance += output.Value
	}

	fmt.Printf("Balance of %s: %d\n", address, balance)
}

func (cli *Cli) Send(from, to string, amount int) {
	chain := blockchain.ContinueBlockChain(from)
	defer chain.Database.Close()

	tx := blockchain.NewTransaction(from, to, amount, chain)
	chain.AddBlock([]*blockchain.Transaction{tx})
	fmt.Println("Success!")
}
func (cli *Cli) Run() {
	cli.ValidateArgs()

	getBalanceCmd := flag.NewFlagSet("balance", flag.ExitOnError)
	createCmd := flag.NewFlagSet("create", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printBlockCmd := flag.NewFlagSet("print", flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address", "", "The address of the account")
	createAddress := createCmd.String("address", "", "The address of the creator")
	sendFrom := sendCmd.String("from", "", "Source wallet")
	sendTo := sendCmd.String("to", "", "Destination wallet")
	sendAmount := sendCmd.Int("amount", 0, "Tx amount")

	switch os.Args[1] {
	case "balance":
		err := getBalanceCmd.Parse(os.Args[2:])
		blockchain.Handle(err)
	case "create":
		err := createCmd.Parse(os.Args[2:])
		blockchain.Handle(err)
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		blockchain.Handle(err)
	case "print":
		err := printBlockCmd.Parse(os.Args[2:])
		blockchain.Handle(err)
	default:
		cli.printUsageAndExit()
	}
	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			cli.printUsageAndExit()
		}
		cli.Balance(*getBalanceAddress)
	}
	if createCmd.Parsed() {
		if *createAddress == "" {
			cli.printUsageAndExit()
		}
		cli.CreateBlockChain(*createAddress)
	}
	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount == 0 {
			cli.printUsageAndExit()
		}
		cli.Send(*sendFrom, *sendTo, *sendAmount)
	}
	if printBlockCmd.Parsed() {
		cli.PrintChain()
	}
}

func (cli *Cli) printUsageAndExit() {
	cli.PrintUsage()
	runtime.Goexit()
}
