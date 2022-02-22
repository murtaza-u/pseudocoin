package cli

import (
	"errors"

	"github.com/murtaza-udaipurwala/pseudocoin/core"
)

type balance struct {
	Address string `json:"address"`
	Balance uint   `json:"balance"`
}

type balanceParams struct {
	Address string `json:"address,omitempty"`
}

func (cli *CLI) getBalance(addr string) (interface{}, error) {
	if !core.ValidateAddress(addr) {
		return nil, errors.New("Invalid address")
	}

	var balance balance

	err := cli.rpcCall("RPC.GetBalance", &balanceParams{Address: addr}, &balance)
	if err != nil {
		return nil, err
	}

	return balance, nil
}
