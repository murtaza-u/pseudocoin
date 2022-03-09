package jsonrpc

import (
	"encoding/hex"
	"net/http"

	"github.com/murtaza-udaipurwala/pseudocoin/core"
)

type PrevTXs struct {
	PrevTXs map[string][]byte `json:"prevTXs"`
}

func (rpc *RPC) GetPrevTXs(r *http.Request, args *struct{ TX []byte }, resp *PrevTXs) error {
	tx, err := core.DeserializeTX(args.TX)
	if err != nil {
		return err
	}

	prevTXs := make(map[string][]byte)

	for _, in := range tx.Inputs {
		prevTX, err := blockchain.FindTXByID(in.TxID)
		if err != nil {
			return err
		}

		serial, err := prevTX.Serialize()
		if err != nil {
			return err
		}

		prevTXs[hex.EncodeToString(prevTX.ID)] = serial
	}

	resp.PrevTXs = prevTXs
	return nil
}
