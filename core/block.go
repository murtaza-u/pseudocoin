package core

type Block struct {
	timestamp     int64
	nonce         uint64
	prevBlockHash []byte
	hash          []byte
	transactions  []*Transaction
}
