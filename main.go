package main

import (
	"fmt"
	"mrnagy.com/owlchain/blockchain"
)

func main() {
	c := blockchain.Genesis()
	fmt.Printf("Chain has %d blocks\n", len(c.Blocks))

	c.AddBlock("First Block")
	fmt.Printf("Chain has %d blocks\n", len(c.Blocks))

	c.AddBlock("Second Blockerado")
	fmt.Printf("Chain has %d blocks\n", len(c.Blocks))

	c.AddBlock("Third Blockage")
	fmt.Printf("Chain has %d blocks\n", len(c.Blocks))

	for _, b := range c.Blocks {
		fmt.Printf("Hash: %x\n", b.Hash)
		fmt.Printf("Data: %s\n", b.Data)
		fmt.Printf("PrevHash: %x\n", b.PrevHash)
		fmt.Printf("Nonce: %5d  Valid: %t\n\n", b.Nonce, blockchain.NewProof(b).Validate())
	}
}
