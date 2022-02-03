package core_test

import (
	"testing"
	"time"

	"github.com/murtaza-udaipurwala/core/core"
)

func TestSerializeDeserialize(t *testing.T) {
	b := core.Block{
		Timestamp:     time.Now().Unix(),
		Nonce:         10,
		PrevBlockHash: []byte{},
		Hash:          []byte{},
		Transactions:  []*core.Transaction{},
	}

	result, err := b.Serialize()
	if err != nil {
		t.Error(err)
	}

	b, err = core.DeserializeBlock(result)
	if err != nil {
		t.Error(err)
	}

	if b.Nonce != 10 {
		t.Errorf("Block deserialization failed")
	}
}

func TestHashTXs(t *testing.T) {
	tx1 := core.Transaction{ID: []byte("1")}
	tx2 := core.Transaction{ID: []byte("2")}
	tx3 := core.Transaction{ID: []byte("3")}

	b := core.Block{Transactions: []*core.Transaction{&tx1, &tx2, &tx3}}

	hash, err := b.HashTXs()
	if err != nil {
		t.Error(err)
	}

	t.Logf("Hash: %s", hash)
}
