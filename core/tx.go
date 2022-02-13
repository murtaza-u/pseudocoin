package core

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
)

type Transaction struct {
	ID      []byte
	Inputs  []TXInput
	Outputs []TXOutput
}

func (tx Transaction) Hash() ([]byte, error) {
	serialData, err := tx.Serialize()
	if err != nil {
		return []byte{}, err
	}

	hash := sha256.Sum256(serialData)
	return hash[:], nil
}

func (tx Transaction) Serialize() ([]byte, error) {
	var buff bytes.Buffer
	encoder := gob.NewEncoder(&buff)
	err := encoder.Encode(tx)
	return buff.Bytes(), err
}

func DeserializeTX(data []byte) (Transaction, error) {
	tx := Transaction{}
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&tx)
	return tx, err
}

const subsidy = 100

func NewCBTX(address, data string) (Transaction, error) {
	if len(data) == 0 {
		randData := make([]byte, 20)
		_, err := rand.Read(randData)
		if err != nil {
			return Transaction{}, err
		}

		data = fmt.Sprintf("%x", randData)
	}

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
