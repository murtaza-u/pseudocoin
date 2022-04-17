package web

import (
	"github.com/murtaza-udaipurwala/pseudocoin/core"
	"github.com/murtaza-udaipurwala/pseudocoin/jsonrpc"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) CreateWallet() (*Wallet, error) {
	w, err := core.NewWallet()
	if err != nil {
		return nil, err
	}

	priv, err := w.EncodePrivKeys()
	if err != nil {
		return nil, err
	}

	pub := w.EncodePubKeys()

	return &Wallet{pub, priv}, nil
}

func (s *Service) GetBalance(addr string) (*jsonrpc.Balance, error) {
	var balance jsonrpc.Balance

	err := jsonrpc.RPCCall(
		"RPC.GetBalance",
		&balanceParams{Address: addr},
		&balance,
	)
	if err != nil {
		return nil, err
	}

	return &balance, nil
}

func (s *Service) Send(r *Send, sender string) (*jsonrpc.Send, error) {
	var newTX jsonrpc.NewTX
	err := jsonrpc.RPCCall("RPC.NewTX", &txParams{
		Sender:       sender,
		Receiver:     r.RecvAddr,
		SenderPubKey: r.SenderPub,
		Amount:       r.Amount,
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

	w := core.Wallet{}
	err = w.DecodePrivKeys(r.SenderPriv)
	if err != nil {
		return nil, err
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

	return &send, nil
}

func (s *Service) GetBlocks(q *BlockQuery) (*jsonrpc.Blocks, error) {
	var blocks jsonrpc.Blocks
	err := jsonrpc.RPCCall(
		"RPC.GetBlocks",
		q,
		&blocks,
	)
	if err != nil {
		return nil, err
	}

	return &blocks, nil
}

func (s *Service) GetAddress(pub string) (string, error) {
	w := core.Wallet{}
	err := w.DecodePubKeys(pub)
	if err != nil {
		return "", ErrInvalidPubKey
	}

	addr, err := w.GetAddress()
	if err != nil {
		return "", err
	}

	return addr, nil
}
