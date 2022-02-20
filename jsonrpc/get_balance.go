package jsonrpc

import (
	"errors"
	"net/http"

	"github.com/mr-tron/base58"
	"github.com/murtaza-udaipurwala/pseudocoin/core"
)

func (rpc *RPC) GetBalance(r *http.Request, args *struct {
	Address string
}, resp *struct {
	Address string `json:"address,omitempty"`
	Balance uint   `json:"balance,omitempty"`
}) error {
	if !core.ValidateAddress(args.Address) {
		return errors.New("Invalid Address")
	}

	var balance uint

	pubKeyHash, err := base58.Decode(args.Address)
	if err != nil {
		return err
	}

	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	UTXOs, err := UTXOSet.FindUTXOs(pubKeyHash)
	if err != nil {
		return err
	}

	for _, out := range UTXOs {
		balance += out.Value
	}

	resp.Address = args.Address
	resp.Balance = balance
	return nil
}
