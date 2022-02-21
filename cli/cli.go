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

func NewCLI() CLI {
	return CLI{}
}

func (cli *CLI) ValidateArgs() error {
	if len(os.Args) < 2 {
		return errors.New("No arguments provided")
	}

	return nil
}

func (cli *CLI) Run() (interface{}, error) {
	err := cli.ValidateArgs()
	if err != nil {
		return nil, err
	}

	return nil, nil
}
