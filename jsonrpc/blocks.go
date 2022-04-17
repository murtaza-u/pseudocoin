package jsonrpc

import (
	"errors"
	"net/http"

	"github.com/murtaza-udaipurwala/pseudocoin/core"
)

type Blocks struct {
	Blocks []core.Block `json:"blocks"`
	Count  uint         `json:"count"`
}

func getBlockHeight(bc *core.Blockchain) (uint, error) {
	i := bc.Iterator()
	var height uint

	for {
		b, err := i.Next()
		if err != nil {
			return 0, err
		}

		if b == nil {
			break
		}

		height++
	}

	return height, nil
}

func (rpc *RPC) GetBlocks(r *http.Request, args *struct{ MaxHT, MinHT uint }, resp *Blocks) error {
	if args.MaxHT < args.MinHT {
		return errors.New("max height cannot be less than min height")
	}

	bc, err := getBlockchain()
	if err != nil {
		return err
	}
	defer bc.DB.Close()

	count, err := getBlockHeight(bc)
	if err != nil {
		return err
	}

	resp.Count = count

	i := bc.Iterator()
	if err != nil {
		return err
	}

	var idx uint

	for {
		b, err := i.Next()
		if err != nil {
			return err
		}

		if b == nil || args.MaxHT < idx {
			break
		}

		if idx < args.MinHT {
			idx++
			continue
		}

		resp.Blocks = append(resp.Blocks, *b)
		idx++
	}

	return nil
}
