package jsonrpc

import (
	"net/http"

	"github.com/murtaza-udaipurwala/pseudocoin/core"
)

type blocks struct {
	Blocks []core.Block `json:"blocks"`
}

func getBlockHeight() (uint, error) {
	i := blockchain.Iterator()
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

func (rpc *RPC) GetBlocks(r *http.Request, args *struct{ Height uint }, resp *blocks) error {
	i := blockchain.Iterator()
	height, err := getBlockHeight()
	if err != nil {
		return err
	}

	for {
		b, err := i.Next()
		if err != nil {
			return err
		}

		if b == nil || args.Height >= height {
			break
		}

		resp.Blocks = append(resp.Blocks, *b)
		height--
	}

	return nil
}
