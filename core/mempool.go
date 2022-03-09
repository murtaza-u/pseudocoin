package core

import (
	"bytes"
	"errors"

	"github.com/boltdb/bolt"
)

const pool = "pool"

func (bc *Blockchain) Add(tx Transaction) error {
	valid, err := bc.VerifyTX(tx)
	if err != nil {
		return err
	}

	if !valid {
		return errors.New("invalid TX")
	}

	serial, err := tx.Serialize()
	if err != nil {
		return err
	}

	err = bc.DB.Update(func(t *bolt.Tx) error {
		b := t.Bucket([]byte(pool))
		if b == nil {
			b, err = t.CreateBucket([]byte(pool))
			if err != nil {
				return err
			}

			return b.Put(tx.ID, serial)
		}

		c := b.Cursor()

		var referenced bool

		for k, v := c.First(); k != nil; k, v = c.Next() {
			ftx, err := DeserializeTX(v)
			if err != nil {
				return err
			}

		Work:
			for _, in := range tx.Inputs {
				for _, fin := range ftx.Inputs {
					if bytes.Compare(in.TxID, fin.TxID) != 0 {
						continue
					}

					if in.Out == fin.Out {
						referenced = true
						break Work
					}
				}
			}
		}

		if referenced {
			return errors.New("Invalid TX")
		}

		return b.Put(tx.ID, serial)
	})

	return err
}
