package core

import (
	"bytes"

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
