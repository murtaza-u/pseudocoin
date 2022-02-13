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

const subsidy = 100

func NewCBTX(address, data string) (Transaction, error) {
	in := TXInput{
		TxID:      []byte{},
		Out:       -1,
		PublicKey: []byte(data),
		Signature: nil,
	}

	out := NewTXOutput(subsidy, address)
	tx := Transaction{
		ID:      []byte{},
		Inputs:  []TXInput{in},
		Outputs: []TXOutput{out},
	}

	txHash, err := tx.Hash()
	if err != nil {
		return Transaction{}, err
	}

	tx.ID = txHash
	return tx, nil
}
