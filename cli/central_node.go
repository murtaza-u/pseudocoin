package cli

import (
	"github.com/murtaza-udaipurwala/pseudocoin/core"
	"github.com/murtaza-udaipurwala/pseudocoin/jsonrpc"
)

func (cli *CLI) startCentralNode(dbFile string) (interface{}, error) {
	chain, err := core.NewBlockchain(dbFile)
	if err != nil {
		return nil, err
	}

	cli.Blockchain = chain
	cli.UTXOSet = core.UTXOSet{Blockchain: &chain}
	cli.UTXOSet.Reindex()
	cli.Blockchain.DB.Close()

	return nil, jsonrpc.InitRPCServer(dbFile)
}
