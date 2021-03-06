package core_test

import (
	"testing"

	"github.com/murtaza-udaipurwala/pseudocoin/core"
)

func TestMerkleTree(t *testing.T) {
	data := [][]byte{
		[]byte("node1"),
		[]byte("node2"),
		[]byte("node3"),
	}

	tree := core.NewMerkleTree(data)
	t.Logf("Root node: %x", tree.Root.Data)
}
