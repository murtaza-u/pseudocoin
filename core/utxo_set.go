package core

import (
	"encoding/hex"
	"errors"

	"github.com/boltdb/bolt"
)

type UTXOSet struct {
	Blockchain *Blockchain
}

const utxoBucket = "utxo"

func (u UTXOSet) Reindex() error {
	db := u.Blockchain.DB
	bucketName := []byte(utxoBucket)

	err := db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket(bucketName)
		if err != nil && err != bolt.ErrBucketNotFound {
			return err
		}

		_, err = tx.CreateBucket(bucketName)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	UTXO, err := u.Blockchain.FindUXTOs()
	if err != nil {
		return err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)

		for txID, outs := range UTXO {
			key, err := hex.DecodeString(txID)
			if err != nil {
				return err
			}

			serial, err := outs.Serialize()
			if err != nil {
				return err
			}

			err = b.Put(key, serial)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

func (u UTXOSet) Update(block *Block) error {
	db := u.Blockchain.DB

	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))

		for _, tx := range block.Transactions {
			if tx.IsCoinbase() == false {
				for _, vin := range tx.Inputs {
					updatedOuts := TXOutputs{}
					outsBytes := b.Get(vin.TxID)
					outs, err := DeserializeOutputs(outsBytes)

					if err != nil {
						return err
					}

					for outIdx, out := range outs.Outputs {
						if outIdx != vin.Out {
							updatedOuts.Outputs = append(updatedOuts.Outputs, out)
						}
					}

					if len(updatedOuts.Outputs) == 0 {
						err := b.Delete(vin.TxID)
						if err != nil {
							return err
						}
					} else {
						serial, err := updatedOuts.Serialize()
						if err != nil {
							return err
						}

						err = b.Put(vin.TxID, serial)
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

func (u *UTXOSet) FindUTXOs(pubkeyHash []byte) ([]TXOutput, error) {
	var UTXOs []TXOutput
	db := u.Blockchain.DB

	err := db.View(func(t *bolt.Tx) error {
		b := t.Bucket([]byte(utxoBucket))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			outs, err := DeserializeOutputs(v)
			if err != nil {
				return err
			}

			for _, out := range outs.Outputs {
				if out.IsLockedWith(pubkeyHash) {
					UTXOs = append(UTXOs, out)
				}
			}
		}

		return nil
	})

	return UTXOs, err
}

func (u *UTXOSet) FindSpendableOutputs(pubkeyHash []byte, amount uint) (uint, map[string][]int, error) {
	if amount == 0 {
		return 0, nil, errors.New("invalid amount")
	}

	var accumulated uint
	spendableOutputs := make(map[string][]int)

	err := u.Blockchain.DB.View(func(t *bolt.Tx) error {
		b := t.Bucket([]byte(utxoBucket))
		c := b.Cursor()

	Accumulate:
		for k, v := c.First(); k != nil; k, v = c.Next() {
			txID := hex.EncodeToString(k)

			outs, err := DeserializeOutputs(v)
			if err != nil {
				return err
			}

			for outIDX, out := range outs.Outputs {
				if out.IsLockedWith(pubkeyHash) {
					spendableOutputs[txID] = append(spendableOutputs[txID], outIDX)
					accumulated += out.Value

					if accumulated >= amount {
						break Accumulate
					}
				}
			}
		}

		return nil
	})

	return accumulated, spendableOutputs, err
}
