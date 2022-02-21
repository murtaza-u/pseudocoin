package cli

import (
	"bytes"
	"errors"
	"net/http"

	"github.com/gorilla/rpc/json"
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

	msg, err := json.EncodeClientRequest("RPC.GetBalance", &balanceParams{Address: addr})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(msg))
	req.Header.Set("Content-Type", "application/json")

	client := new(http.Client)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var balance balance
	err = json.DecodeClientResponse(resp.Body, &balance)
	if err != nil {
		return nil, err
	}

	return balance, nil
}
