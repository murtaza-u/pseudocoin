package jsonrpc

import (
	"errors"
	"net/http"

	"github.com/mr-tron/base58"
	"github.com/murtaza-udaipurwala/pseudocoin/core"
)

type MyTX struct {
	TxID     []byte `json:"txid"`
	Amount   uint   `json:"amount"`
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Msg      string `json:"msg"`
}

type MyTXs struct {
	TXs []MyTX `json:"txs"`
}

func (rpc *RPC) GetMyTXs(r *http.Request, args *struct{ Address string }, resp *MyTXs) error {
	if !core.ValidateAddress(args.Address) {
		return errors.New("Invalid Address")
	}

	payload, err := base58.Decode(args.Address)
	if err != nil {
		return err
	}

	pubKeyHash := payload[1 : len(payload)-4]

	bc, err := getBlockchain()
	if err != nil {
		return err
	}
	defer bc.DB.Close()

	i := bc.Iterator()
	for {
		b, err := i.Next()
		if err != nil {
			return err
		}

		if b == nil {
			break
		}

		for _, tx := range b.Transactions {
			if tx.IsCoinbase() {
				continue
			}

			out := tx.Outputs[0]

			if out.IsLockedWith(pubKeyHash) {
				in := tx.Inputs[0]

				w := core.Wallet{}
				w.PubKey = in.PublicKey
				addr, err := w.GetAddress()
				if err != nil {
					return err
				}

				myTX := MyTX{
					Amount:   out.Value,
					Sender:   addr,
					Receiver: args.Address,
					Msg:      tx.Msg,
					TxID:     tx.ID,
				}

				resp.TXs = append(resp.TXs, myTX)
				continue
			}

			uses, err := tx.Inputs[0].UsesKey(pubKeyHash)
			if err != nil {
				return err
			}

			if !uses {
				continue
			}

			var amount uint
			for _, in := range tx.Inputs {
				tx, err := bc.FindTXByID(in.TxID)
				if err != nil {
					return err
				}

				amount += tx.Outputs[in.Out].Value
			}

			versionPayload := append([]byte{core.Version}, out.PubkeyHash...)
			checksum := core.Checksum(versionPayload)
			fullPayload := append(versionPayload, checksum...)
			addr := base58.Encode(fullPayload)

			myTX := MyTX{
				Amount:   amount,
				Sender:   args.Address,
				Receiver: addr,
				Msg:      tx.Msg,
				TxID:     tx.ID,
			}

			resp.TXs = append(resp.TXs, myTX)
		}
	}

	return nil
}
