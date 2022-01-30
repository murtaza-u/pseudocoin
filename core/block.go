package core

import (
	"bytes"
	"encoding/gob"
	"log"
)

type Block struct {
	Timestamp     int64
	Nonce         uint64
	PrevBlockHash []byte
	Hash          []byte
	Transactions  []*Transaction
}

func (b Block) serialize() []byte {
	var buff bytes.Buffer
	encoder := gob.NewEncoder(&buff)
	if err := encoder.Encode(b); err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

func deserializeBlock(encData []byte) Block {
	block := Block{}
	decoder := gob.NewDecoder(bytes.NewReader(encData))
	if err := decoder.Decode(&block); err != nil {
		log.Panic(err)
	}
	return block
}
