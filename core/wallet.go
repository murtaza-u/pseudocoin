package core

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"

	"github.com/mr-tron/base58"
	"golang.org/x/crypto/ripemd160"
)

type Wallet struct {
	PrivKey ecdsa.PrivateKey
	PubKey  []byte
}

const (
	addressChecksumLen = 4
	version            = byte(0x00)
)

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

func HashPubKey(pub []byte) ([]byte, error) {
	pubSHA256 := sha256.Sum256(pub)

	RIPEMD160 := ripemd160.New()
	_, err := RIPEMD160.Write(pubSHA256[:])
	if err != nil {
		return nil, err
	}

	pubRIPEMD160 := RIPEMD160.Sum(nil)
	return pubRIPEMD160, nil
}

func Checksum(pubKeyHash []byte) []byte {
	firstSHA256 := sha256.Sum256(pubKeyHash)
	secondSHA256 := sha256.Sum256(firstSHA256[:])

	// checksum is only first 4 bytes of the resulting hash
	return secondSHA256[:addressChecksumLen]
}

func (w Wallet) GetAddress() (string, error) {
	pubKeyHash, err := HashPubKey(w.PubKey)
	if err != nil {
		return "", err
	}

	versionPayload := append([]byte{version}, pubKeyHash...)
	checksum := Checksum(versionPayload)
	fullPayload := append(versionPayload, checksum...)

	return base58.Encode(fullPayload), nil
}

func ValidateAddress(address string) bool {
	payload, err := base58.Decode(address)
	if err != nil {
		return false
	}

	checksum := payload[len(payload)-addressChecksumLen:]
	version := payload[0]
	pubKeyHash := payload[1 : len(payload)-addressChecksumLen]

	targetChecksum := Checksum(append([]byte{version}, pubKeyHash...))
	return bytes.Compare(checksum, targetChecksum) == 0
}
