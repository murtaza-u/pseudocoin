package jsonrpc

import (
	"errors"
	"net/http"

	"github.com/mr-tron/base58"
	"github.com/murtaza-udaipurwala/pseudocoin/core"
)

type Balance struct {
	Address string `json:"address"`
	Balance uint   `json:"balance"`
}

func (rpc *RPC) GetBalance(r *http.Request, args *struct{ Address string }, resp *Balance) error {
	if !core.ValidateAddress(args.Address) {
		return errors.New("Invalid Address")
	}

	payload, err := base58.Decode(args.Address)
	if err != nil {
		return err
	}

	pubKeyHash := payload[1 : len(payload)-4]
	outs, err := utxoset.FindUTXOs(pubKeyHash)
	if err != nil {
		return err
	}

	var balance uint
	for _, out := range outs {
		balance += out.Value
	}

	resp.Address = args.Address
	resp.Balance = balance
	return nil
}
