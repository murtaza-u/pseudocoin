package core_test

import (
	"bytes"
	"testing"

	"github.com/murtaza-udaipurwala/pseudocoin/core"
)

func TestTXHash(t *testing.T) {
	tx := core.Transaction{}
	hash, err := tx.Hash()

	if err != nil {
		t.Error(err)
	}

	if len(hash) == 0 {
		t.Errorf("Invalid Hash")
	}
}

func TestSerializeDeserializeTX(t *testing.T) {
	tx := core.Transaction{ID: []byte("123")}
	data, err := tx.Serialize()
	if err != nil {
		t.Error(err)
	}

	tx, err = core.DeserializeTX(data)
	if err != nil {
		t.Error(err)
	}

	if bytes.Compare(tx.ID, []byte("123")) != 0 {
		t.Errorf("Transaction deserialization failed")
	}
}
