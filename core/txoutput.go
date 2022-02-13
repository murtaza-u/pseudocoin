package core

import (
	"bytes"
	"encoding/gob"

	"github.com/mr-tron/base58"
)

type TXOutput struct {
	Value      uint
	PubkeyHash []byte
}

// sign the output
func (out *TXOutput) Lock(address string) error {
	pubKeyHash, err := base58.Decode(address)
	if err != nil {
		return err
	}

	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-addressChecksumLen]
	out.PubkeyHash = pubKeyHash
	return nil
}

// check if the output can be unlocked
func (out TXOutput) IsLockedWith(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubkeyHash, pubKeyHash) == 0
}

// create a new TXOutput
func NewTXOutput(value uint, address string) TXOutput {
	out := TXOutput{
		Value:      value,
		PubkeyHash: nil,
	}

	out.Lock(address)
	return out
}

type TXOutputs struct {
	Outputs []TXOutput
}

func (outs TXOutputs) Serialize() ([]byte, error) {
	var buff bytes.Buffer
	encoder := gob.NewEncoder(&buff)
	if err := encoder.Encode(outs); err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}

func DeserializeOutputs(encData []byte) (TXOutputs, error) {
	decoder := gob.NewDecoder(bytes.NewReader(encData))
	outputs := TXOutputs{}
	if err := decoder.Decode(&outputs); err != nil {
		return TXOutputs{}, err
	}

	return outputs, nil
}
