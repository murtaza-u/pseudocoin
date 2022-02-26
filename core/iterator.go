package core

import "github.com/boltdb/bolt"

type Iterator struct {
	currentBlockHash []byte
	db               *bolt.DB
}

func (i *Iterator) Next() (*Block, error) {
	if len(i.currentBlockHash) == 0 {
		return nil, nil
	}

	var block Block

	err := i.db.View(func(t *bolt.Tx) error {
		b := t.Bucket([]byte(blocksBucket))
		data := b.Get(i.currentBlockHash)

		var err error
		block, err = DeserializeBlock(data)

		return err
	})

	if err != nil {
		return nil, err
	}

	i.currentBlockHash = block.PrevBlockHash

	return &block, err
}
