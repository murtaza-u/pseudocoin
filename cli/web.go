package cli

import "github.com/murtaza-udaipurwala/pseudocoin/web"

func (cli *CLI) web() {
	web.Init()
}
