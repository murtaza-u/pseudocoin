package cli

import (
	"errors"
	"os"

	"github.com/murtaza-udaipurwala/pseudocoin/core"
)

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

func (cli *CLI) ValidateArgs() error {
	if len(os.Args) < 2 {
		return errors.New("No arguments provided")
	}

	return nil
}
