package core

type Transaction struct {
	ID      []byte
	Inputs  []TXInput
	Outputs []TXOutput
}
