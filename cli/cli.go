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
	walletCMDName := walletCMD.String("name", "", "Give a name to the wallet")

	addressCMD := flag.NewFlagSet("getaddress", flag.ExitOnError)
	addressCMDFile := addressCMD.String("i", "", "Specify public key file")
	addressCMDPubKey := addressCMD.String("pubkey", "", "Pass public key")

	flag.Parse()
	config := Config{}
	config.Load(*configFile)

	switch os.Args[1] {
	case "wallet":
		err = walletCMD.Parse(os.Args[2:])

	case "getaddress":
		err = addressCMD.Parse(os.Args[2:])
	}

	if err != nil {
		return nil, err
	}

	if walletCMD.Parsed() {
		if *walletCMDCreate {
			return cli.CreateWallet(*walletCMDName)
		}
	}

	if addressCMD.Parsed() {
		if len(*addressCMDFile) != 0 {
			return cli.GetAddress(*addressCMDFile)
		}

		if len(*addressCMDPubKey) != 0 {
			return cli.GetAddress(*addressCMDPubKey)
		}
	}

	return nil, errors.New("Invalid arguments")
}
