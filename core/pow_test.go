package core_test

import (
	"testing"
	"time"

	"github.com/murtaza-udaipurwala/core/core"
)

var pow = core.NewPoW(&core.Block{
	Timestamp:     time.Now().Unix(),
	PrevBlockHash: []byte("0"),
	Transactions: []*core.Transaction{
		{ID: []byte("1")},
		{ID: []byte("2")},
		{ID: []byte("3")},
	},
})
var hash []byte
var nonce uint64

func TestPrepareData(t *testing.T) {
	_, err := pow.PrepareData(10)
	if err != nil {
		t.Error(err)
	}
}

func TestRun(t *testing.T) {
	var err error
	hash, nonce, err = pow.Run()

	if err != nil {
		t.Error(err)
	}

	if len(hash) == 0 {
		t.Errorf("Invalid hash\n")
	}

	if nonce == 0 {
		t.Errorf("Invalid nonce")
	}

	t.Logf("Nonce: %v", nonce)
	pow.Block.Nonce = nonce
}

func TestValidate(t *testing.T) {
	isValid, err := pow.Validate()
	if err != nil {
		t.Error(err)
	}

	if !isValid {
		t.Errorf("Work done is not valid\n")
	}
}
