package cli

import (
	"errors"

	"github.com/murtaza-udaipurwala/pseudocoin/core"
)

type blockchain struct {
	Address string `json:"address"`
	DBFile  string `json:"dbfile"`
}

func (cli *CLI) createBlockchain(dbFile, addr string) (interface{}, error) {
	if !core.ValidateAddress(addr) {
		return nil, errors.New("invalid address provided")
	}

	chain, err := core.CreateBlockchain(addr, dbFile)
	if err != nil {
		return nil, err
	}

	cli.Blockchain = chain
	cli.UTXOSet = core.UTXOSet{Blockchain: &chain}

	// reindex UTXO set
	cli.UTXOSet.Reindex()

	return blockchain{
		Address: addr,
		DBFile:  dbFile,
	}, nil
}
