package core_test

import (
	"testing"

	"github.com/murtaza-udaipurwala/core/core"
)

var wallet core.Wallet

func TestNewWallet(t *testing.T) {
	var err error
	wallet, err = core.NewWallet()
	if err != nil {
		t.Fatal(err)
	}
}

func TestHashPubKey(t *testing.T) {
	hashed, err := core.HashPubKey(wallet.PubKey)
	if err != nil {
		t.Fatal(err)
	}

	if len(hashed) == 0 {
		t.Fatal("invalid pub key hash")
	}
}

func TestChecksum(t *testing.T) {
	checksum := core.Checksum(wallet.PubKey)
	if len(checksum) != 4 {
		t.Fatal("invalid checksum")
	}
}

func TestGetAddress(t *testing.T) {
	address, err := wallet.GetAddress()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(address)
	if !core.ValidateAddress(address) {
		t.Fatal("invalid address")
	}
}
