package core

import (
	"errors"

	"github.com/boltdb/bolt"
)

const (
	dbFile              = "blockchain.db"
	blocksBucket        = "blocks"
	genesisCoinbaseData = "I use Arch BTW"
)

type Blockchain struct {
	Tip []byte
	DB  *bolt.DB
}

func CreateBlockchain(address string) (Blockchain, error) {
	if DBExists(dbFile) {
		return Blockchain{}, errors.New("Blockchain already exists")
	}

	cbtx, err := NewCBTX(address, genesisCoinbaseData)
	if err != nil {
		return Blockchain{}, err
	}

	genesis, err := NewGenesisBlock(cbtx)
	if err != nil {
		return Blockchain{}, err
	}

	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		return Blockchain{}, err
	}

	err = db.Update(func(t *bolt.Tx) error {
		b, err := t.CreateBucket([]byte(blocksBucket))
		if err != nil {
			return err
		}

		serial, err := genesis.Serialize()
		if err != nil {
			return err
		}

		err = b.Put(genesis.Hash, serial)
		if err != nil {
			return err
		}

		err = b.Put([]byte("l"), genesis.Hash)
		if err != nil {
			return err
		}

		return nil
	})

	return Blockchain{
		Tip: genesis.Hash,
		DB:  db,
	}, err
}

func NewBlockchain() (Blockchain, error) {
	if !DBExists(dbFile) {
		return Blockchain{}, errors.New("Blockchain does not exists. Create one first")
	}

	var tip []byte

	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		return Blockchain{}, err
	}

	db.View(func(t *bolt.Tx) error {
		b := t.Bucket([]byte(blocksBucket))
		tip = b.Get([]byte("l"))
		return nil
	})

	return Blockchain{
		Tip: tip,
		DB:  db,
	}, nil
}

func (bc *Blockchain) iterator() *Iterator {
	return &Iterator{
		currentBlockHash: bc.Tip,
		db: bc.DB,
	}
}
