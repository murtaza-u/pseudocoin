package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/murtaza-udaipurwala/pseudocoin/cli"
)

func main() {
	cli := cli.NewCLI()

	res, err := cli.Run()
	if err != nil {
		log.Println(err)
		return
	}

	b, err := json.MarshalIndent(res, "", "    ")
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(string(b))
}
