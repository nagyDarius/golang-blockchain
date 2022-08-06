package main

import (
	"mrnagy.com/owlchain/cli"
	"os"
)

func main() {
	defer os.Exit(0)
	c := cli.Cli{}
	c.Run()
}
