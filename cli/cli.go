package cli

import "github.com/murtaza-udaipurwala/pseudocoin/core"

type CLI struct {
	Blockchain core.Blockchain
	UTXOSet    core.UTXOSet
}
