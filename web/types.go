package web

import "errors"

type balanceParams struct {
	Address string `json:"address,omitempty"`
}

type Wallet struct {
	PubKey  string
	PrivKey string
}

type Send struct {
	RecvAddr   string `json:"recv_addr"`
	SenderPub  string `json:"sender_pub"`
	SenderPriv string `json:"sender_priv"`
	Amount     uint   `json:"amount"`
}

type txParams struct {
	Sender       string `json:"sender"`
	Receiver     string `json:"receiver"`
	SenderPubKey string `json:"sender_pub_key"`
	Amount       uint   `json:"amount"`
}

type blockParams struct {
	Height uint `json:"height"`
}

var ErrInvalidPubKey = errors.New("invalid public key")