package core

import (
	"encoding/hex"

	"github.com/boltdb/bolt"
)

type UTXOSet struct {
	Blockchain *Blockchain
}

const utxoBucket = "utxo"

func (u *UTXOSet) Reindex() error {
	db := u.Blockchain.DB

	utxos, err := u.Blockchain.findUXTOs()
	if err != nil {
		return err
	}

	err = db.Update(func(t *bolt.Tx) error {
		err := t.DeleteBucket([]byte(utxoBucket))
		if err != nil && err != bolt.ErrBucketNotFound {
			return err
		}

		b, err := t.CreateBucket([]byte(utxoBucket))
		if err != nil {
			return err
		}

		for txid, outs := range utxos {
			key, err := hex.DecodeString(txid)
			if err != nil {
				return err
			}

			serializedOuts, err := outs.Serialize()
			if err != nil {
				return err
			}

			b.Put(key, serializedOuts)
		}

		return nil
	})

	return err
}
