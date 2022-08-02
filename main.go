package main

import (
	"github.com/dgraph-io/badger"
	"mrnagy.com/owlchain/blockchain"
	"mrnagy.com/owlchain/cli"
	"os"
)

func main() {
	defer os.Exit(0)
	chain := blockchain.NewChain()
	defer func(Database *badger.DB) {
		err := Database.Close()
		if err != nil {
			blockchain.Handle(err)
		}
	}(chain.Database)

	c := cli.Cli{Chain: chain}
	c.Run()
}
