package cli

import (
	"encoding/json"
	"io/ioutil"
)

type config struct {
	DB      string `json:"db"`
	Address string `json:"node_address"`
}

func (c *config) load(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, c)
}
