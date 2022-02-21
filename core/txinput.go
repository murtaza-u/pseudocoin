package core

import "bytes"

type TXInput struct {
	TxID      []byte `json:"txid"`
	Out       int    `json:"out"`
	PublicKey []byte `json:"public_key"`
	Signature []byte `json:"signature"`
}

// check whether the address initiated the transaction
func (in TXInput) UsesKey(pubKeyHash []byte) (bool, error) {
	lockingHash, err := HashPubKey(in.PublicKey)
	if err != nil {
		return false, err
	}

	return bytes.Compare(lockingHash, pubKeyHash) == 0, nil
}
