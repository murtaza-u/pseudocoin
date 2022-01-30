package core_test

import (
	"testing"

	"github.com/murtaza-udaipurwala/core/core"
)

func TestTXHash(t *testing.T) {
	tx := core.Transaction{}
	hash, err := tx.Hash()

	if err != nil {
		t.Error(err)
	}

	if len(hash) == 0 {
		t.Errorf("Invalid Hash\n")
	}
}
