package jsonrpc

import (
	"encoding/hex"
	"net/http"

	"github.com/murtaza-udaipurwala/pseudocoin/core"
)

type prevTXs struct {
	PrevTXs map[string]core.Transaction `json:"prevTXs"`
}

func (rpc *RPC) GetPrevTXs(r *http.Request, args *struct{ TX []byte }, resp *prevTXs) error {
	tx, err := core.DeserializeTX(args.TX)
	if err != nil {
		return err
	}

	prevTXs := make(map[string]core.Transaction)

	for _, in := range tx.Inputs {
		prevTX, err := blockchain.FindTXByID(in.TxID)
		if err != nil {
			return err
		}

		prevTXs[hex.EncodeToString(in.TxID)] = prevTX
	}

	resp.PrevTXs = prevTXs
	return nil
}
