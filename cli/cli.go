package cli

import (
	"errors"
	"flag"
	"os"
	"path"

	"github.com/murtaza-udaipurwala/pseudocoin/core"
)

type CLI struct {
	Blockchain core.Blockchain
	UTXOSet    core.UTXOSet
	Config     config
}

const URL = "http://localhost:5000/rpc"

func NewCLI() CLI {
	return CLI{}
}

func (cli *CLI) validateArgs() error {
	if len(os.Args) < 2 {
		return errors.New("No arguments provided")
	}

	return nil
}

func (cli *CLI) Run() (interface{}, error) {
	err := cli.validateArgs()
	if err != nil {
		return nil, err
	}

	createHome()
	home, err := getHome()
	if err != nil {
		return nil, err
	}

	configFile := flag.String("config", path.Join(home, "config.json"), "Path to the config file")

	walletCMD := flag.NewFlagSet("wallet", flag.ExitOnError)
	walletCMDCreate := walletCMD.Bool("create", false, "Create a new wallet")
	walletCMDName := walletCMD.String("name", "", "Give a name to the wallet")

	addressCMD := flag.NewFlagSet("getaddress", flag.ExitOnError)
	addressCMDFile := addressCMD.String("i", "", "Specify public key file")
	addressCMDPubKey := addressCMD.String("pubkey", "", "Pass public key")

	blockchainCMD := flag.NewFlagSet("blockchain", flag.ExitOnError)
	blockchainCMDCreate := blockchainCMD.String("create", "", "Create a new blockchain")

	centralNodeCMD := flag.NewFlagSet("centralnode", flag.ExitOnError)
	centralNodeCMDStart := centralNodeCMD.String("start", "", "Path to blockchain DB")

	getBalanceCMD := flag.NewFlagSet("getbalance", flag.ExitOnError)
	getBalanceCMDAddr := getBalanceCMD.String("addr", "", "Pseudocoin address")

	flag.Parse()
	cli.Config.load(*configFile)

	switch os.Args[1] {
	case "wallet":
		err = walletCMD.Parse(os.Args[2:])

	case "getaddress":
		err = addressCMD.Parse(os.Args[2:])

	case "blockchain":
		err = blockchainCMD.Parse(os.Args[2:])

	case "centralnode":
		err = centralNodeCMD.Parse(os.Args[2:])

	case "getbalance":
		err = getBalanceCMD.Parse(os.Args[2:])
	}

	if err != nil {
		return nil, err
	}

	if walletCMD.Parsed() {
		if *walletCMDCreate {
			return cli.createWallet(*walletCMDName)
		}
	}

	if addressCMD.Parsed() {
		if len(*addressCMDFile) != 0 {
			return cli.getAddress(*addressCMDFile)
		}

		if len(*addressCMDPubKey) != 0 {
			return cli.getAddress(*addressCMDPubKey)
		}
	}

	if blockchainCMD.Parsed() {
		if len(*blockchainCMDCreate) != 0 {
			return cli.createBlockchain(*blockchainCMDCreate)
		}
	}

	if centralNodeCMD.Parsed() {
		if len(*centralNodeCMDStart) != 0 {
			return cli.startCentralNode(*centralNodeCMDStart)
		}
	}

	if getBalanceCMD.Parsed() {
		if len(*getBalanceCMDAddr) != 0 {
			return cli.getBalance(*getBalanceCMDAddr)
		}
	}

	return nil, errors.New("Invalid arguments")
}
