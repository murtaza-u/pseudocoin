package miner

import (
	"encoding/hex"
	"log"
	"time"

	"github.com/murtaza-udaipurwala/pseudocoin/core"
	"github.com/murtaza-udaipurwala/pseudocoin/jsonrpc"
)

func getTXs() ([]*core.Transaction, error) {
	var pool jsonrpc.Pool

	err := jsonrpc.RPCCall(
		"RPC.GetPool",
		nil,
		&pool,
	)
	if err != nil {
		return nil, err
	}

	var txs []*core.Transaction

	for _, tx := range pool.TXs {
		d, err := hex.DecodeString(tx)
		if err != nil {
			return nil, err
		}

		t, err := core.DeserializeTX(d)
		if err != nil {
			return nil, err
		}

		txs = append(txs, &t)
	}

	return txs, nil
}

func getPrevHash() ([]byte, error) {
	var prev jsonrpc.PrevHash

	err := jsonrpc.RPCCall(
		"RPC.GetPrevBlockHash",
		nil,
		&prev,
	)
	if err != nil {
		return nil, err
	}

	return prev.Hash, nil
}

func mine(addr string, txs []*core.Transaction) (*core.Block, error) {
	prevHash, err := getPrevHash()
	if err != nil {
		log.Panic(err)
		return nil, err
	}

	cbtx, err := core.NewCBTX(addr, "")
	if err != nil {
		log.Panic(err)
		return nil, err
	}

	txs = append(txs, &cbtx)

	b, err := core.NewBlock(txs, prevHash)
	if err != nil {
		log.Panic(err)
		return nil, err
	}

	return &b, nil
}

func report(b *core.Block) error {
	serial, err := b.Serialize()
	if err != nil {
		return err
	}

	var r jsonrpc.Report
	err = jsonrpc.RPCCall(
		"RPC.ReportBlock",
		&jsonrpc.ReportParams{Block: serial},
		&r,
	)
	if err != nil {
		return err
	}

	log.Println(r.Msg)
	return nil
}

func Start(addr string) {
	start := func() error {
		txs, err := getTXs()
		if err != nil {
			return err
		}

		if len(txs) == 0 {
			return nil
		}

		b, err := mine(addr, txs)
		if err != nil {
			return err
		}

		err = report(b)
		return err
	}

	for {
		err := start()
		if err != nil {
			log.Println("error:", err)
		}

		time.Sleep(time.Second * 10)
	}
}

// func Start(addr string) {
// 	dur := time.Now().Add(time.Second * 5)

// 	start := func() error {
// 		txs, err := getTXs()
// 		if err != nil {
// 			return err
// 		}

// 		if len(txs) == 0 {
// 			dur = time.Now().Add(time.Second * 5)
// 			return nil
// 		}

// 		if len(txs) < 2 && time.Now().Sub(dur) < 0 {
// 			return nil
// 		}

// 		dur = time.Now().Add(time.Second * 5)

// 		b, err := mine(addr, txs)
// 		if err != nil {
// 			return err
// 		}

// 		err = report(b)
// 		return err
// 	}

// 	for {
// 		err := start()
// 		if err != nil {
// 			log.Println("error:", err)
// 		}

// 		time.Sleep(time.Second * 10)
// 	}
// }
