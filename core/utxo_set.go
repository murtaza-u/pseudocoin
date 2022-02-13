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

func (u *UTXOSet) update(block Block) error {
	db := u.Blockchain.DB

	err := db.Update(func(t *bolt.Tx) error {
		b := t.Bucket([]byte(utxoBucket))

		for _, tx := range block.Transactions {
			if !tx.IsCoinbase() {
				for _, in := range tx.Inputs {
					updatedOutputs := TXOutputs{}
					outs, err := DeserializeOutputs(b.Get(in.TxID))
					if err != nil {
						return err
					}

					for outIDX, out := range outs.Outputs {
						if outIDX == in.Out {
							updatedOutputs.Outputs = append(updatedOutputs.Outputs, out)
						}
					}

					if len(updatedOutputs.Outputs) == 0 {
						err = b.Delete(in.TxID)
						if err != nil {
							return err
						}
					} else {
						serial, err := updatedOutputs.Serialize()
						if err != nil {
							return err
						}

						err = b.Put(in.TxID, serial)
						if err != nil {
							return err
						}
					}
				}
			}

			newOutputs := TXOutputs{}
			for _, out := range tx.Outputs {
				newOutputs.Outputs = append(newOutputs.Outputs, out)
			}

			serial, err := newOutputs.Serialize()
			if err != nil {
				return err
			}

			err = b.Put(tx.ID, serial)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}
