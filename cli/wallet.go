package cli

import "github.com/murtaza-udaipurwala/pseudocoin/core"

type wallet struct {
	Address string `json:"address"`
	Msg     string `json:"msg"`
}

func (cli *CLI) CreateWallet(name string) (interface{}, error) {
	var err error

	if len(name) == 0 {
		name, err = randomString()
		if err != nil {
			return nil, err
		}
	}

	w, err := core.NewWallet()
	if err != nil {
		return nil, err
	}

	addr, err := w.GetAddress()
	if err != nil {
		return nil, err
	}

	// write private keys
	priv, err := w.EncodePrivKeys()
	if err != nil {
		return nil, err
	}

	if err := write(name, []byte(priv)); err != nil {
		return nil, err
	}

	// write public keys
	if err := write(name+".pub", []byte(w.EncodePubKeys())); err != nil {
		return nil, err
	}

	return wallet{
		Address: addr,
		Msg:     "Wallet created successfully",
	}, nil
}
