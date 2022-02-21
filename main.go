package main

import "github.com/murtaza-udaipurwala/pseudocoin/cli"

func main() {
	cli := cli.NewCLI()
	cli.Print(cli.Run())
}
