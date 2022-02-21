package cli

import (
	"errors"
	"flag"
	"os"

	"github.com/murtaza-udaipurwala/pseudocoin/core"
)

type CLI struct {
	Blockchain core.Blockchain
	UTXOSet    core.UTXOSet
}

func NewCLI() CLI {
	return CLI{}
}

func (cli *CLI) ValidateArgs() error {
	if len(os.Args) < 2 {
		return errors.New("No arguments provided")
	}

	return nil
}

func (cli *CLI) Run() (interface{}, error) {
	err := cli.ValidateArgs()
	if err != nil {
		return nil, err
	}

	configFile := flag.String("config", ".config.json", "Path to the config file")

	walletCMD := flag.NewFlagSet("wallet", flag.ExitOnError)
	walletCMDCreate := walletCMD.Bool("create", false, "Create a new wallet")

	flag.Parse()
	config := Config{}
	config.Load(*configFile)

	switch os.Args[1] {
	case "wallet":
		err = walletCMD.Parse(os.Args[2:])
	}

	if err != nil {
		return nil, err
	}

	if walletCMD.Parsed() {
		if *walletCMDCreate {

		}
	}

	return nil, errors.New("Invalid arguments")
}
