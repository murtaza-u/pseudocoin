package cli

import (
	"errors"
	"flag"
	"os"
	"path"

	"github.com/murtaza-udaipurwala/pseudocoin/core"
	"github.com/murtaza-udaipurwala/pseudocoin/miner"
)

type CLI struct {
	Blockchain core.Blockchain
	UTXOSet    core.UTXOSet
	Config     config
}

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

	sendCMD := flag.NewFlagSet("send", flag.ExitOnError)
	sendCMDSender := sendCMD.String("sender", "", "Sender's address")
	sendCMDRecv := sendCMD.String("recv", "", "receiver's address")
	sendCMDAmount := sendCMD.Uint("amount", 0, "amount")
	sendCMDPriv := sendCMD.String("priv", "", "sender's private key")
	sendCMDPub := sendCMD.String("pub", "", "sender's public key")

	poolCMD := flag.NewFlagSet("pool", flag.ExitOnError)

	mineCMD := flag.NewFlagSet("mine", flag.ExitOnError)
	mineCMDAddr := mineCMD.String("addr", "", "miner's address")

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

	case "send":
		err = sendCMD.Parse(os.Args[2:])

	case "pool":
		err = poolCMD.Parse(os.Args[2:])

	case "mine":
		err = mineCMD.Parse(os.Args[2:])
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

	if sendCMD.Parsed() {
		if len(*sendCMDSender) != 0 && len(*sendCMDRecv) != 0 && len(*sendCMDPriv) != 0 && len(*sendCMDPub) != 0 && *sendCMDAmount != 0 {
			return cli.send(*sendCMDRecv, *sendCMDSender, *sendCMDPriv, *sendCMDPub, *sendCMDAmount)
		}
	}

	if poolCMD.Parsed() {
		return cli.getPool()
	}

	if mineCMD.Parsed() {
		if len(*mineCMDAddr) != 0 {
			miner.Start(*mineCMDAddr)
			os.Exit(1)
		}
	}

	return nil, errors.New("Invalid arguments")
}
