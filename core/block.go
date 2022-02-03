package core

import (
	"bytes"
	"encoding/gob"
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
		serialData, err := tx.serialize()
		if err != nil {
			return nil, err
		}
		txs = append(txs, serialData)
	}

	tree := NewMerkleTree(txs)
	return tree.Root.Data, nil
}
