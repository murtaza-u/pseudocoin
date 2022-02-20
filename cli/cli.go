package cli

import "github.com/murtaza-udaipurwala/pseudocoin/core"

type CLI struct {
	Blockchain core.Blockchain
	UTXOSet    core.UTXOSet
}

func NewCLI(bc core.Blockchain, UTXOSet core.UTXOSet) CLI {
	return CLI{
		Blockchain: bc,
		UTXOSet:    UTXOSet,
	}
}
