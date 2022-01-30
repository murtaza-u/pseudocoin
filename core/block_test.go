package core_test

import (
	"testing"
	"time"

	"github.com/murtaza-udaipurwala/core/core"
)

func TestSerializeDeserialize(t *testing.T) {
	b := core.Block {
		Timestamp: time.Now().Unix(),
		Nonce: 10,
		PrevBlockHash: []byte{},
		Hash: []byte{},
		Transactions: []*core.Transaction{},
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
