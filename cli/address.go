package cli

import (
	"io/ioutil"

	"github.com/murtaza-udaipurwala/pseudocoin/core"
)

type address struct {
	Address string `json:"address"`
}

func (cli *CLI) getAddress(arg string) (interface{}, error) {
	var pubkey []byte
	var err error

	if fileExists(arg) {
		pubkey, err = ioutil.ReadFile(arg)
		if err != nil {
			return nil, err
		}
	} else {
		pubkey = []byte(arg)
	}

	w := core.Wallet{}
	if err = w.DecodePubKeys(string(pubkey)); err != nil {
		return nil, err
	}

	addr, err := w.GetAddress()
	if err != nil {
		return nil, err
	}

	return address{
		Address: addr,
	}, nil
}
