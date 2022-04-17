package cli

import (
	"errors"
	"io/ioutil"

	"github.com/murtaza-udaipurwala/pseudocoin/core"
	"github.com/murtaza-udaipurwala/pseudocoin/jsonrpc"
)

type txParams struct {
	Sender       string `json:"sender"`
	Receiver     string `json:"receiver"`
	SenderPubKey string `json:"sender_pub_key"`
	Amount       uint   `json:"amount"`
	Msg          string `json:"msg"`
}

func (cli *CLI) send(receiver, sender, senderPriv, senderPub, msg string, amount uint) (interface{}, error) {
	if !core.ValidateAddress(receiver) || !core.ValidateAddress(sender) {
		return nil, errors.New("Invalid address")
	}

	var privKey, pubKey string
	var err error

	// private key
	if fileExists(senderPriv) {
		p, err := ioutil.ReadFile(senderPriv)
		if err != nil {
			return nil, err
		}

		privKey = string(p)
	} else {
		privKey = senderPriv
	}

	// public key
	if fileExists(senderPub) {
		p, err := ioutil.ReadFile(senderPub)
		if err != nil {
			return nil, err
		}

		pubKey = string(p)
	} else {
		pubKey = senderPub
	}

	w := core.Wallet{}
	w.DecodePrivKeys(privKey)

	var newTX jsonrpc.NewTX
	err = jsonrpc.RPCCall("RPC.NewTX", &txParams{
		Sender:       sender,
		Receiver:     receiver,
		SenderPubKey: pubKey,
		Amount:       amount,
		Msg:          msg,
	}, &newTX)
	if err != nil {
		return nil, err
	}

	transaction, err := core.DeserializeTX(newTX.TX)
	if err != nil {
		return nil, err
	}

	var prevTXs jsonrpc.PrevTXs
	err = jsonrpc.RPCCall("RPC.GetPrevTXs", &jsonrpc.NewTX{
		TX: newTX.TX,
	}, &prevTXs)
	if err != nil {
		return nil, err
	}

	ptx := make(map[string]core.Transaction)
	for txID, serial := range prevTXs.PrevTXs {
		tx, err := core.DeserializeTX(serial)
		if err != nil {
			return nil, err
		}

		ptx[txID] = tx
	}

	err = transaction.Sign(w.PrivKey, ptx)
	if err != nil {
		return nil, err
	}

	serial, err := transaction.Serialize()
	if err != nil {
		return nil, err
	}

	var send jsonrpc.Send
	err = jsonrpc.RPCCall("RPC.Send", &jsonrpc.NewTX{
		TX: serial,
	}, &send)
	if err != nil {
		return nil, err
	}

	return send, nil
}
