package cli

import (
	"errors"
	"io/ioutil"

	"github.com/murtaza-udaipurwala/pseudocoin/core"
)

type tx struct {
	TX []byte `json:"tx"`
}

type txParams struct {
	Sender       string `json:"sender"`
	Receiver     string `json:"receiver"`
	SenderPubKey string `json:"sender_pub_key"`
	Amount       uint   `json:"amount"`
}

type prevTXs struct {
	PrevTXs map[string]core.Transaction `json:"prevTXs"`
}

type send struct {
	Msg string `json:"msg"`
}

func (cli *CLI) send(receiver, sender, senderPriv, senderPub string, amount uint) (interface{}, error) {
	if !core.ValidateAddress(receiver) || !core.ValidateAddress(sender) {
		return nil, errors.New("Invalid address")
	}

	var privKey, pubKey []byte
	var err error

	// private key
	if fileExists(senderPriv) {
		privKey, err = ioutil.ReadFile(senderPriv)
		if err != nil {
			return nil, err
		}
	} else {
		privKey = []byte(senderPriv)
	}

	// public key
	if fileExists(senderPub) {
		pubKey, err = ioutil.ReadFile(senderPub)
		if err != nil {
			return nil, err
		}
	} else {
		pubKey = []byte(senderPub)
	}

	w := core.Wallet{}
	w.DecodePrivKeys(string(privKey))

	var newTx tx
	err = cli.rpcCall("RPC.NewTX", &txParams{
		Sender:       sender,
		Receiver:     receiver,
		SenderPubKey: string(pubKey),
		Amount:       amount,
	}, &newTx)
	if err != nil {
		return nil, err
	}

	var prevTXs prevTXs
	err = cli.rpcCall("RPC.GetPrevTXs", &tx{
		TX: newTx.TX,
	}, &prevTXs)
	if err != nil {
		return nil, err
	}

	transaction, err := core.DeserializeTX(newTx.TX)
	if err != nil {
		return nil, err
	}

	err = transaction.Sign(w.PrivKey, prevTXs.PrevTXs)
	if err != nil {
		return nil, err
	}

	serial, err := transaction.Serialize()
	if err != nil {
		return nil, err
	}

	var send send
	err = cli.rpcCall("RPC.Send", &tx{
		TX: serial,
	}, &send)
	if err != nil {
		return nil, err
	}

	return send, nil
}
