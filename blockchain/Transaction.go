package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

type Transaction struct {
	ID      []byte
	Inputs  []TxInput
	Outputs []TxOutput
}

type TxInput struct {
	ID  []byte
	Out int
	Sig string
}

type TxOutput struct {
	Value     int
	PublicKey string
}

func CoinbaseTx(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Coins to %s", to)
	}

	txin := TxInput{[]byte{}, -1, data}
	txout := TxOutput{100, to}

	tx := Transaction{nil, []TxInput{txin}, []TxOutput{txout}}
	tx.setID()

	return &tx
}

func (tx *Transaction) setID() {
	var encoded bytes.Buffer

	err := gob.NewEncoder(&encoded).Encode(tx)
	Handle(err)

	hash := sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].ID) == 0 && tx.Inputs[0].Out == -1
}

func (i *TxInput) CanUnlock(data string) bool {
	return i.Sig == data
}

func (o *TxOutput) CanBeUnlocked(data string) bool {
	return o.PublicKey == data
}

func NewTransaction(from, to string, amount int, c *Chain) *Transaction {
	var in []TxInput
	var out []TxOutput

	acc, validOutputs := c.FindSpendableOutputs(from, amount)

	if acc < amount {
		log.Panic("Error: Not enough funds")
	}

	for id, outs := range validOutputs {
		txId, err := hex.DecodeString(id)
		Handle(err)
		for _, out := range outs {
			input := TxInput{txId, out, from}
			in = append(in, input)
		}
	}

	out = append(out, TxOutput{amount, to})

	if acc > amount {
		out = append(out, TxOutput{acc - amount, from})
	}

	tx := Transaction{nil, in, out}
	tx.setID()
	return &tx
}
