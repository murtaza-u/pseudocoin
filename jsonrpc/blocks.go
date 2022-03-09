package jsonrpc

import (
	"net/http"

	"github.com/murtaza-udaipurwala/pseudocoin/core"
)

type Blocks struct {
	Blocks []core.Block `json:"blocks"`
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

func (rpc *RPC) GetBlocks(r *http.Request, args *struct{ Height uint }, resp *Blocks) error {
	bc, err := getBlockchain()
	if err != nil {
		return err
	}

	i := bc.Iterator()
	height, err := getBlockHeight(bc)
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
