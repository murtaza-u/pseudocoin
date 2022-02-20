package core

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"

	"github.com/boltdb/bolt"
)

const (
	blocksBucket        = "blocks"
	genesisCoinbaseData = "I use Arch BTW"
)

type Blockchain struct {
	Tip []byte
	DB  *bolt.DB
}

func CreateBlockchain(address, dbFile string) (Blockchain, error) {
	if DBExists(dbFile) {
		return Blockchain{}, errors.New("Blockchain already exists")
	}

	cbtx, err := NewCBTX(address, genesisCoinbaseData)
	if err != nil {
		return Blockchain{}, err
	}

	genesis, err := NewGenesisBlock(cbtx)
	if err != nil {
		return Blockchain{}, err
	}

	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		return Blockchain{}, err
	}

	err = db.Update(func(t *bolt.Tx) error {
		b, err := t.CreateBucket([]byte(blocksBucket))
		if err != nil {
			return err
		}

		serial, err := genesis.Serialize()
		if err != nil {
			return err
		}

		err = b.Put(genesis.Hash, serial)
		if err != nil {
			return err
		}

		err = b.Put([]byte("l"), genesis.Hash)
		if err != nil {
			return err
		}

		return nil
	})

	return Blockchain{
		Tip: genesis.Hash,
		DB:  db,
	}, err
}

func NewBlockchain(dbFile string) (Blockchain, error) {
	if !DBExists(dbFile) {
		return Blockchain{}, errors.New("Blockchain does not exists. Create one first")
	}

	var tip []byte

	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		return Blockchain{}, err
	}

	db.View(func(t *bolt.Tx) error {
		b := t.Bucket([]byte(blocksBucket))
		tip = b.Get([]byte("l"))
		return nil
	})

	return Blockchain{
		Tip: tip,
		DB:  db,
	}, nil
}

func (bc *Blockchain) iterator() *Iterator {
	return &Iterator{
		currentBlockHash: bc.Tip,
		db:               bc.DB,
	}
}

func (bc *Blockchain) FindTXByID(ID []byte) (Transaction, error) {
	i := bc.iterator()

	for {
		b, err := i.Next()
		if err != nil {
			return Transaction{}, err
		}

		if b == nil {
			break
		}

		for _, tx := range b.Transactions {
			if bytes.Compare(tx.ID, ID) == 0 {
				return *tx, nil
			}
		}
	}

	return Transaction{}, errors.New("Transaction not found")
}

func (bc *Blockchain) VerifyTX(tx Transaction) (bool, error) {
	if tx.IsCoinbase() {
		return true, nil
	}

	prevTXs := make(map[string]Transaction)

	for _, in := range tx.Inputs {
		prevTX, err := bc.FindTXByID(in.TxID)
		if err != nil {
			return false, err
		}

		prevTXs[hex.EncodeToString(in.TxID)] = prevTX
	}

	return tx.Verify(prevTXs)
}

func (bc *Blockchain) SignTX(tx Transaction, privKey ecdsa.PrivateKey) error {
	if tx.IsCoinbase() {
		return nil
	}

	prevTXs := make(map[string]Transaction)

	for _, in := range tx.Inputs {
		prevTX, err := bc.FindTXByID(in.TxID)
		if err != nil {
			return err
		}

		prevTXs[hex.EncodeToString(in.TxID)] = prevTX
	}

	return tx.Sign(privKey, prevTXs)
}

func (bc *Blockchain) MineBlock(txs []*Transaction) (Block, error) {
	for _, tx := range txs {
		isVerified, err := bc.VerifyTX(*tx)
		if err != nil {
			return Block{}, errors.New("Failed to verify transaction")
		}

		if !isVerified {
			return Block{}, errors.New("One or more invalid transaction(s)")
		}
	}

	newBlock, err := NewBlock(txs, bc.Tip)
	if err != nil {
		return Block{}, err
	}

	err = bc.DB.Update(func(t *bolt.Tx) error {
		b := t.Bucket([]byte(blocksBucket))
		err := b.Put([]byte("l"), newBlock.Hash)
		if err != nil {
			return err
		}

		serialBlock, err := newBlock.Serialize()
		if err != nil {
			return err
		}

		return b.Put(newBlock.Hash, serialBlock)
	})
	if err != nil {
		return Block{}, err
	}

	bc.Tip = newBlock.Hash
	return newBlock, nil
}

// scan the entire blockchain and find all UTXOs
func (chain *Blockchain) FindUXTOs() (map[string]TXOutputs, error) {
	UTXOs := make(map[string]TXOutputs)
	spentUTXOs := make(map[string][]int)
	i := chain.iterator()

	for {
		block, err := i.Next()
		if err != nil {
			return UTXOs, err
		}

		if block == nil {
			break
		}

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

			if tx.IsCoinbase() != true {
				for _, in := range tx.Inputs {
					inID := hex.EncodeToString(in.TxID)
					spentUTXOs[inID] = append(spentUTXOs[inID], in.Out)
				}
			}

		Output:
			for outIDX, out := range tx.Outputs {
				if spentUTXOs[txID] != nil {
					for _, spentOut := range spentUTXOs[txID] {
						if spentOut == outIDX {
							continue Output
						}
					}
				}

				outs := UTXOs[txID]
				outs.Outputs = append(outs.Outputs, out)
				UTXOs[txID] = outs
			}
		}
	}

	return UTXOs, nil
}
