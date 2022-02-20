package core_test

import (
	"testing"

	"github.com/murtaza-udaipurwala/core/core"
)

func TestNext(t *testing.T) {
	i := core.Iterator{}
	_, err := i.Next()
	if err != nil {
		t.Error(err)
	}
}
