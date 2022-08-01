package blockchain

import (
	"github.com/dgraph-io/badger"
)

type Chain struct {
	LastHash []byte
	Database *badger.DB
}

func (c *Chain) AddBlock(data string) {
	prevBlock := c.Blocks[len(c.Blocks)-1]
	newBlock := CreateBlock(data, prevBlock.Hash)
	c.Blocks = append(c.Blocks, newBlock)
}

func Genesis() *Chain {
	return &Chain{[]*Block{CreateBlock("OwlChain", []byte{})}}
}
