package core_test

import (
	"os"
	"testing"

	"github.com/murtaza-udaipurwala/pseudocoin/core"
)

const (
	blockchainExistsErr = "Blockchain already exists"
	db                  = "test.db"
)

var bc core.Blockchain

func TestCreateBlockchain(t *testing.T) {
	w, err := core.NewWallet()
	if err != nil {
		t.Error(err)
	}

	address, err := w.GetAddress()
	if err != nil {
		t.Error(err)
	}

	bc, err = core.CreateBlockchain(address, db)
	if err != nil {
		if err.Error() == blockchainExistsErr {
			t.SkipNow()
		}
		t.Error(err)
	}

	t.Logf("Tip: %x\n", bc.Tip)
}

func TestFindUTXOs(t *testing.T) {
	_, err := bc.FindUXTOs()
	if err != nil {
		t.Error(err)
	}
}

func TestCleanUp(t *testing.T) {
	os.Remove(db)
}
