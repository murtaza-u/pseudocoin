package cli

import (
	"errors"

	"github.com/murtaza-udaipurwala/pseudocoin/core"
)

type blockchain struct {
	Address string `json:"address"`
	DBFile  string `json:"dbfile"`
}

func (cli *CLI) CreateBlockchain(dbFile string) (interface{}, error) {
	if len(cli.Config.Address) == 0 {
		return nil, errors.New("Please provide node's address in config.json")
	}

	if !core.ValidateAddress(cli.Config.Address) {
		return nil, errors.New("Invalid address provided in config.json")
	}

	chain, err := core.CreateBlockchain(cli.Config.Address, dbFile)
	if err != nil {
		return nil, err
	}

	cli.Blockchain = chain
	cli.UTXOSet = core.UTXOSet{Blockchain: &chain}

	// reindex UTXO set
	cli.UTXOSet.Reindex()

	return blockchain{
		Address: cli.Config.Address,
		DBFile:  dbFile,
	}, nil
}
