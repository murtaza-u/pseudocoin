package core_test

import (
	"testing"

	"github.com/murtaza-udaipurwala/pseudocoin/core"
)

func TestIntToBytes(t *testing.T) {
	result, err := core.IntToBytes(10)
	if err != nil {
		t.Error(err)
	}

	if len(result) == 0 {
		t.Errorf("Invalid result\n")
	}
}
