package jsonrpc

import (
	"net/http"

	"github.com/murtaza-udaipurwala/pseudocoin/core"
)

type send struct {
	Msg string `json:"msg"`
}

func (rpc *RPC) Send(r *http.Request, args *struct{ TX []byte }, resp *send) error {
	tx, err := core.DeserializeTX(args.TX)
	if err != nil {
		return err
	}

	err = mempool.Add(tx)
	if err != nil {
		return err
	}

	resp.Msg = "Block added to the mempool"
	return nil
}
