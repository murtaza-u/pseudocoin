package core

type Block struct {
	Timestamp     int64
	Nonce         uint64
	PrevBlockHash []byte
	Hash          []byte
	Transactions  []*Transaction
}
