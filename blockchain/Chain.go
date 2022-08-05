package blockchain

import (
	"encoding/hex"
	"fmt"
	"github.com/dgraph-io/badger"
	"os"
	"runtime"
)

const (
	dbPath      = "./tmp/owl-blocks"
	dbFile      = "./tmp/owl-blocks/MANIFEST"
	genesisData = "First transaction from genesis"
)

type Chain struct {
	LastHash []byte
	Database *badger.DB
}

func InitBlockChain(address string) *Chain {
	if dbExists() {
		fmt.Println("Blockchain already exists")
		runtime.Goexit()
	}

	db, err := badger.Open(badger.DefaultOptions(dbPath))
	Handle(err)

	var lastHash []byte
	err = db.Update(func(txn *badger.Txn) error {
		lastHash, err = createNewChain(address, txn)
		return err
	})
	Handle(err)

	return &Chain{
		LastHash: lastHash,
		Database: db,
	}
}

func ContinueBlockChain(address string) *Chain {
	if !dbExists() {
		fmt.Println("Blockchain does not exist. Create one!")
		runtime.Goexit()
	}

	db, err := badger.Open(badger.DefaultOptions(dbPath))
	Handle(err)

	var lastHash []byte

	err = db.Update(func(txn *badger.Txn) error {
		lastHashDB, err := txn.Get([]byte("lh"))
		Handle(err)
		err = lastHashDB.Value(func(val []byte) error {
			lastHash = val
			return nil
		})
		return err
	})

	return &Chain{
		LastHash: lastHash,
		Database: db,
	}
}

func dbExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}
	return true
}

func createNewChain(address string, txn *badger.Txn) ([]byte, error) {
	fmt.Println("No existing blockchain found")
	genesis := genesis(CoinbaseTx(address, genesisData))
	fmt.Println("genesis proved")
	err := txn.Set(genesis.Hash, genesis.Serialize())
	Handle(err)
	err = txn.Set([]byte("lh"), genesis.Hash)

	return genesis.Hash, err
}

func genesis(coinbase *Transaction) *Block {
	return CreateBlock([]*Transaction{coinbase}, []byte{})
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

func (c *Chain) FindUnspentTransactions(address string) []*Transaction {
	var unspentTxs []*Transaction

	spentTxOs := make(map[string][]int)

	iter := c.Iterator()

	for {
		block := iter.Next()

		for _, tx := range block.Transactions {
			txId := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Outputs {
				if spentTxOs[txId] != nil {
					for _, spentOut := range spentTxOs[txId] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}
				if out.CanBeUnlocked(address) {
					unspentTxs = append(unspentTxs, tx)
				}
			}
			if tx.IsCoinbase() == false {
				for _, in := range tx.Inputs {
					if in.CanUnlock(address) {
						inTxID := hex.EncodeToString(in.ID)
						spentTxOs[inTxID] = append(spentTxOs[inTxID], in.Out)
					}
				}
			}
		}
		if len(block.PrevHash) == 0 {
			break
		}
	}
	return unspentTxs
}
