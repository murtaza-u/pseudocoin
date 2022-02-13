package core

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
)

type Wallet struct {
	PrivKey ecdsa.PrivateKey
	PubKey  []byte
}

func newKeyPair() (ecdsa.PrivateKey, []byte, error) {
	curve := elliptic.P256()
	priv, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return ecdsa.PrivateKey{}, nil, err
	}

	pub := append(priv.PublicKey.X.Bytes(), priv.PublicKey.Y.Bytes()...)
	return *priv, pub, nil
}

func NewWallet() (Wallet, error) {
	priv, pub, err := newKeyPair()
	if err != nil {
		return Wallet{}, err
	}

	wallet := Wallet{
		PrivKey: priv,
		PubKey:  pub,
	}

	return wallet, nil
}
