package jsonrpc

import (
	"net/http"

	"github.com/murtaza-udaipurwala/pseudocoin/core"
)

type Send struct {
	Msg string `json:"msg"`
}

func (rpc *RPC) Send(r *http.Request, args *struct{ TX []byte }, resp *Send) error {
	bc, err := getBlockchain()
	if err != nil {
		return err
	}
	defer bc.DB.Close()

	tx, err := core.DeserializeTX(args.TX)
	if err != nil {
		return err
	}

	err = bc.AddToPool(tx)
	if err != nil {
		return err
	}

	resp.Msg = "TX added to the mempool"
	return nil
}
