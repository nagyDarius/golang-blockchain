package blockchain

import "github.com/dgraph-io/badger"

type ChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

func (c *Chain) Iterator() *ChainIterator {
	return &ChainIterator{c.LastHash, c.Database}
}

func (iter *ChainIterator) Next() *Block {
	var block *Block

	err := iter.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.CurrentHash)
		err = item.Value(func(val []byte) error {
			block = Deserialize(val)
			return nil
		})
		return err
	})

	Handle(err)

	iter.CurrentHash = block.PrevHash

	return block
}
