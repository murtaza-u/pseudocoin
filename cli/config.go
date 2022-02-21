package cli

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	DB string `json:"db"`
}

func (c *Config) Load(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, c)
}
