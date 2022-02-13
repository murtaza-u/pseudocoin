package core

import "github.com/mr-tron/base58"

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
