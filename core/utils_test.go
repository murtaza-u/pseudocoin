package core_test

import (
	"testing"

	"github.com/murtaza-udaipurwala/core/core"
)

func TestIntToBytes(t *testing.T) {
	_, err := core.IntToBytes(10)
	if err != nil {
		t.Error(err)
	}
}
