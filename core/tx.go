package core

type Transaction struct {
	id      []byte
	inputs  []TXInput
	outputs []TXOutput
}
