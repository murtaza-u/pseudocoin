package cli

import (
	"encoding/json"
	"fmt"
	"log"
)

type result struct {
	Result interface{} `json:"result,omitempty"`
	Err    string      `json:"err,omitempty"`
}

func (cli *CLI) Print(res interface{}, err error) {
	result := result{Result: res}
	if err != nil {
		result.Err = err.Error()
	}

	data, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(data))
}
