package core

import "github.com/boltdb/bolt"

type Iterator struct {
	currentBlockHash []byte
	db               *bolt.DB
}
