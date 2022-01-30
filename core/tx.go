package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
)

type Transaction struct {
	ID      []byte
	Inputs  []TXInput
	Outputs []TXOutput
}

func (tx Transaction) Hash() ([]byte, error) {
	serialData, err := tx.serialize()
	if err != nil {
		return []byte{}, err
	}

	hash := sha256.Sum256(serialData)
	return hash[:], nil
}

func (tx Transaction) serialize() ([]byte, error) {
	var buff bytes.Buffer
	encoder := gob.NewEncoder(&buff)
	err := encoder.Encode(tx)
	return buff.Bytes(), err
}
