package jsonrpc

import (
	"net/http"

	"github.com/murtaza-udaipurwala/pseudocoin/core"
)

type newTX struct {
	TX []byte `json:"tx"`
}

type newTXArgs struct {
	Sender       string `json:"sender"`
	Receiver     string `json:"receiver"`
	SenderPubKey string `json:"sender_pub_key"`
	Amount       uint   `json:"amount"`
}

func (rpc *RPC) NewTX(r *http.Request, args *newTXArgs, resp *newTX) error {
	w := core.Wallet{}
	err := w.DecodePubKeys(args.SenderPubKey)
	if err != nil {
		return err
	}

	tx, err := core.NewUTXOTransaction(args.Receiver, args.Sender, w.PubKey, args.Amount, &utxoset)
	if err != nil {
		return err
	}

	serial, err := tx.Serialize()
	if err != nil {
		return err
	}

	resp.TX = serial
	return nil
}
