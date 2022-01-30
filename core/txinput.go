package core

type TXInput struct {
	TxID      []byte
	Out       int
	PublicKey []byte
	Signature []byte
}
