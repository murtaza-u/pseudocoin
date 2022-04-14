package cli

import (
	"encoding/hex"

	"github.com/murtaza-udaipurwala/pseudocoin/core"
	"github.com/murtaza-udaipurwala/pseudocoin/jsonrpc"
)

type poolParams struct{}

func (cli *CLI) getPool() (interface{}, error) {
	var pool jsonrpc.Pool

	err := jsonrpc.RPCCall("RPC.GetPool", &poolParams{}, &pool)
	if err != nil {
		return nil, err
	}

	var txs []core.Transaction

	for _, tx := range pool.TXs {
		d, err := hex.DecodeString(tx)
		if err != nil {
			return nil, err
		}

		tx, err := core.DeserializeTX(d)
		if err != nil {
			return nil, err
		}

		txs = append(txs, tx)
	}

	return txs, nil
}
