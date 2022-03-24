package jsonrpc

import "net/http"

type PrevHash struct {
	Hash []byte `json:"prev_hash"`
}

func (rpc *RPC) GetPrevBlockHash(r *http.Request, args *struct{}, resp *PrevHash) error {
	bc, err := getBlockchain()
	if err != nil {
		return err
	}
	defer bc.DB.Close()

	resp.Hash = bc.Tip
	return nil
}
