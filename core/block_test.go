package core_test

import (
	"testing"

	"github.com/murtaza-udaipurwala/pseudocoin/core"
)

var block core.Block

func TestNewBlock(t *testing.T) {
	var err error
	block, err = core.NewBlock(
		[]*core.Transaction{
			{ID: []byte("1")},
			{ID: []byte("2")},
			{ID: []byte("3")},
		}, []byte{},
	)

	if err != nil {
		t.Error(err)
	}
}

func TestSerializeDeserializeBlock(t *testing.T) {
	result, err := block.Serialize()
	if err != nil {
		t.Error(err)
	}

	b, err := core.DeserializeBlock(result)
	if err != nil {
		t.Error(err)
	}

	if b.Nonce != block.Nonce {
		t.Errorf("Block deserialization failed")
	}
}

func TestHashTXs(t *testing.T) {
	hash, err := block.HashTXs()
	if err != nil {
		t.Error(err)
	}

	t.Logf("Hash: %x", hash)
}

func TestNewGenesisBlock(t *testing.T) {
	_, err := core.NewGenesisBlock(core.Transaction{})
	if err != nil {
		t.Error(err)
	}
}
