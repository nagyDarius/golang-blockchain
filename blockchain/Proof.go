package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"
)

const Difficulty = 6

type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

func NewProof(block *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-Difficulty))
	return &ProofOfWork{block, target}
}

func (pow *ProofOfWork) InitData(nonce int) []byte {
	return bytes.Join(
		[][]byte{
			pow.Block.PrevHash,
			pow.Block.HashTransactions(),
			toHex(int64(nonce)),
			toHex(int64(Difficulty)),
		},
		[]byte{},
	)
}

func (pow *ProofOfWork) Run() (int, []byte) {
	var intHash big.Int
	var hash [32]byte

	nonce := 0

	for nonce < math.MaxInt64 {
		data := pow.InitData(nonce)
		hash = sha256.Sum256(data)
		fmt.Printf("Working %d: %X\n", nonce, hash)
		intHash.SetBytes(hash[:])
		if intHash.Cmp(pow.Target) == -1 {
			break
		} else {
			nonce++
		}
	}
	return nonce, hash[:]
}

func (pow *ProofOfWork) Validate() bool {
	var intHash big.Int
	data := pow.InitData(pow.Block.Nonce)

	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:])

	return intHash.Cmp(pow.Target) == -1
}

func toHex(i int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, i)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}
