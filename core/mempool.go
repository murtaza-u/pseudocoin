package core

import (
	"bytes"
	"fmt"
	"sync"
	"time"

	"github.com/boltdb/bolt"
)

type Mempool struct {
	Mutex      *sync.Mutex
	Blockchain *Blockchain
}

const (
	pending = "pending"
	queue   = "queue"
)

func (mem *Mempool) Add(tx Transaction) error {
	mem.Mutex.Lock()

	err := mem.Blockchain.DB.Update(func(t *bolt.Tx) error {
		b, err := t.CreateBucketIfNotExists([]byte(pending))
		if err != nil {
			return err
		}

		serial, err := tx.Serialize()
		if err != nil {
			return err
		}

		return b.Put(tx.ID, serial)
	})

	mem.Mutex.Unlock()

	return err
}

func (mem *Mempool) Queue() error {
	var err error

	for {
		mem.Mutex.Lock()

		err = mem.Blockchain.DB.Update(func(t *bolt.Tx) error {
			pb, err := t.CreateBucketIfNotExists([]byte(pending))
			if err != nil {
				return err
			}

			qb, err := t.CreateBucketIfNotExists([]byte(queue))
			if err != nil {
				return err
			}

			pc := pb.Cursor()
			qc := qb.Cursor()

			for k, v := pc.First(); k != nil; k, v = pc.Next() {
				ptx, err := DeserializeTX(v)
				if err != nil {
					return err
				}

				valid, err := mem.Blockchain.VerifyTX(ptx)
				if err != nil {
					return err
				}

				if !valid {
					if err := pc.Delete(); err != nil {
						return err
					}
					continue
				}

				var referenced bool
				for k, v := qc.First(); k != nil; k, v = qc.Next() {
					qtx, err := DeserializeTX(v)
					if err != nil {
						return err
					}

				Work:
					for _, pin := range ptx.Inputs {
						for _, qin := range qtx.Inputs {
							if bytes.Compare(pin.TxID, qin.TxID) != 0 {
								continue
							}

							if pin.Out == qin.Out {
								referenced = true
								break Work
							}
						}
					}
				}

				if referenced {
					if err := pc.Delete(); err != nil {
						return err
					}
					continue
				}

				if err := pc.Delete(); err != nil {
					return err
				}

				serial, err := ptx.Serialize()
				if err != nil {
					return err
				}

				fmt.Printf("Adding block with ID: %x\n", ptx.ID)
				return qb.Put(ptx.ID, serial)
			}

			return nil
		})

		mem.Mutex.Unlock()

		if err != nil {
			break
		}

		time.Sleep(time.Second * 10)
	}

	return err
}
