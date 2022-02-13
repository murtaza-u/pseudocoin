package core

import (
	"bytes"
	"encoding/gob"
	"time"
)

type Block struct {
	Timestamp     int64
	Nonce         uint64
	PrevBlockHash []byte
	Hash          []byte
	Transactions  []*Transaction
}

func (b Block) Serialize() ([]byte, error) {
	var buff bytes.Buffer
	encoder := gob.NewEncoder(&buff)
	err := encoder.Encode(b)
	return buff.Bytes(), err
}

func DeserializeBlock(encData []byte) (Block, error) {
	block := Block{}
	decoder := gob.NewDecoder(bytes.NewReader(encData))
	err := decoder.Decode(&block)
	return block, err
}

func (b *Block) HashTXs() ([]byte, error) {
	var txs [][]byte

	for _, tx := range b.Transactions {
		serialData, err := tx.Serialize()
		if err != nil {
			return nil, err
		}
		txs = append(txs, serialData)
	}

	tree := NewMerkleTree(txs)
	return tree.Root.Data, nil
}

func NewBlock(txs []*Transaction, prevBlockHash []byte) (Block, error) {
	block := Block{
		Timestamp:     time.Now().Unix(),
		Transactions:  txs,
		PrevBlockHash: prevBlockHash,
		Nonce:         0,
		Hash:          []byte{},
	}

	pow := NewPoW(&block)
	hash, nonce, err := pow.Run()
	if err != nil {
		return Block{}, err
	}

	block.Nonce = nonce
	block.Hash = hash
	return block, nil
}

func NewGenesisBlock(cbtx Transaction) (Block, error) {
	return NewBlock([]*Transaction{&cbtx}, []byte{})
}
