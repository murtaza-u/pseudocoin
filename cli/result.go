package cli

import (
	"encoding/json"
	"fmt"
	"log"
)

type Result struct {
	Result interface{} `json:"result"`
	Err    string      `json:"err"`
}

func (cli *CLI) Print(res interface{}, err error) {
	result := Result{Result: res}
	if err != nil {
		result.Err = err.Error()
	}

	data, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(data))
}
