package jsonrpc

import (
	"encoding/hex"
	"net/http"

	"github.com/boltdb/bolt"
)

type Pool struct {
	TXs []string `json:"txs"`
}

func (rpc *RPC) GetPool(r *http.Request, args *struct{}, resp *Pool) error {
	bc, err := getBlockchain()
	if err != nil {
		return err
	}
	defer bc.DB.Close()

	txs := []string{}

	bc.DB.View(func(t *bolt.Tx) error {
		b := t.Bucket([]byte("pool"))
		if b == nil {
			return nil
		}

		b.ForEach(func(k, v []byte) error {
			d := hex.EncodeToString(v)
			txs = append(txs, d)
			return nil
		})

		return nil
	})

	resp.TXs = txs

	return nil
}
