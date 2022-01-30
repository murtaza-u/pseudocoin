package core

type TXInput struct {
	txID      []byte
	out       int
	publicKey []byte
	signature []byte
}
