package jsonrpc

import (
	"bytes"
	"errors"
	"log"
	"net/http"

	"github.com/boltdb/bolt"
	"github.com/murtaza-udaipurwala/pseudocoin/core"
)

type ReportParams struct {
	Block []byte `json:"block"`
}

type Report struct {
	Msg string `json:"msg"`
}

func (rpc *RPC) ReportBlock(r *http.Request, args *ReportParams, resp *Report) error {
	block, err := core.DeserializeBlock(args.Block)
	if err != nil {
		return err
	}

	log.Printf("hash: %x\n", block.Hash)

	pow := core.NewPoW(&block)
	valid, err := pow.Validate()
	if err != nil {
		return errors.New("failed to validate proof of work")
	}

	if !valid {
		return errors.New("invalid proof of work")
	}

	bc, err := getBlockchain()
	if err != nil {
		return err
	}
	defer bc.DB.Close()

	for _, tx := range block.Transactions {
		err := bc.DB.View(func(t *bolt.Tx) error {
			b := t.Bucket([]byte("blocks"))
			err := b.ForEach(func(k, v []byte) error {
				if bytes.Compare(tx.ID, k) == 0 {
					return errors.New("block already mined")
				}

				return nil
			})

			return err
		})

		if err != nil {
			return err
		}

		valid, err := bc.VerifyTX(*tx)
		if err != nil {
			return errors.New("failed to verify tx")
		}

		if !valid {
			return errors.New("one or more invalid tx")
		}
	}

	serial, err := block.Serialize()
	if err != nil {
		return err
	}

	err = bc.DB.Update(func(t *bolt.Tx) error {
		b := t.Bucket([]byte("blocks"))
		err := b.Put(block.Hash, serial)
		if err != nil {
			return err
		}

		err = b.Put([]byte("l"), block.Hash)
		if err != nil {
			return err
		}

		b = t.Bucket([]byte("pool"))
		for _, tx := range block.Transactions {
			err := b.Delete(tx.ID)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	bc.DB.Close()

	bc, err = getBlockchain()
	if err != nil {
		return err
	}
	defer bc.DB.Close()

	utxoset := core.UTXOSet{Blockchain: bc}
	err = utxoset.Reindex()
	if err != nil {
		return err
	}

	resp.Msg = "block added to the blockchain"
	return nil
}
