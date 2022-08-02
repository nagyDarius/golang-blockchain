package blockchain

import (
	"fmt"
	"github.com/dgraph-io/badger"
)

const dbPath = "./tmp/owl-blocks"

type Chain struct {
	LastHash []byte
	Database *badger.DB
}

func NewChain() *Chain {
	db, err := badger.Open(badger.DefaultOptions(dbPath))
	Handle(err)

	lastHash, err := readOrCreateChain(db)

	Handle(err)

	return &Chain{lastHash, db}
}

func readOrCreateChain(db *badger.DB) ([]byte, error) {
	var lastHash []byte

	err := db.Update(func(txn *badger.Txn) error {
		lastHashDB, err := txn.Get([]byte("lh"))

		if err == badger.ErrKeyNotFound {
			lastHash, err = createNewChain(txn, err)
			return err
		} else if err == nil {
			err = lastHashDB.Value(func(val []byte) error {
				lastHash = val
				return nil
			})
		}
		return err
	})

	return lastHash, err
}

func createNewChain(txn *badger.Txn, err error) ([]byte, error) {
	fmt.Println("No existing blockchain found")
	genesis := genesis()
	fmt.Println("genesis proved")
	err = txn.Set(genesis.Hash, genesis.Serialize())
	Handle(err)
	err = txn.Set([]byte("lh"), genesis.Hash)

	return genesis.Hash, err
}

func genesis() *Block {
	return CreateBlock("OwlChain", []byte{})
}

func (c *Chain) AddBlock(data string) {
	lastHash, err := c.readLastHash()
	Handle(err)

	newBlock := CreateBlock(data, lastHash)

	err = c.addNewBlock(newBlock)
	Handle(err)
}

func (c *Chain) addNewBlock(newBlock *Block) error {
	err := c.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		Handle(err)
		err = txn.Set([]byte("lh"), newBlock.Hash)
		return err
	})
	return err
}

func (c *Chain) readLastHash() ([]byte, error) {
	var lastHash []byte

	err := c.Database.View(func(txn *badger.Txn) error {
		lastHashDB, err := txn.Get([]byte("lh"))
		Handle(err)
		err = lastHashDB.Value(func(val []byte) error {
			lastHash = val
			return nil
		})
		return err
	})
	Handle(err)
	return lastHash, err
}
