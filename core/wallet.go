package core

import "crypto/ecdsa"

type Wallet struct {
	PrivKey ecdsa.PrivateKey
	PubKey  []byte
}
