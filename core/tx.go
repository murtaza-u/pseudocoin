package core

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
)

type Transaction struct {
	ID      []byte     `json:"id"`
	Inputs  []TXInput  `json:"inputs"`
	Outputs []TXOutput `json:"outputs"`
}

func (tx Transaction) Hash() ([]byte, error) {
	serialData, err := tx.Serialize()
	if err != nil {
		return []byte{}, err
	}

	hash := sha256.Sum256(serialData)
	return hash[:], nil
}

func (tx Transaction) Serialize() ([]byte, error) {
	var buff bytes.Buffer
	encoder := gob.NewEncoder(&buff)
	err := encoder.Encode(tx)
	return buff.Bytes(), err
}

func DeserializeTX(data []byte) (Transaction, error) {
	tx := Transaction{}
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&tx)
	return tx, err
}

const subsidy = 100

func NewCBTX(address, data string) (Transaction, error) {
	if len(data) == 0 {
		randData := make([]byte, 20)
		_, err := rand.Read(randData)
		if err != nil {
			return Transaction{}, err
		}

		data = fmt.Sprintf("%x", randData)
	}

	in := TXInput{
		TxID:      []byte{},
		Out:       -1,
		PublicKey: []byte(data),
		Signature: nil,
	}

	out := NewTXOutput(subsidy, address)
	tx := Transaction{
		ID:      []byte{},
		Inputs:  []TXInput{in},
		Outputs: []TXOutput{out},
	}

	txHash, err := tx.Hash()
	if err != nil {
		return Transaction{}, err
	}

	tx.ID = txHash
	return tx, nil
}

func (tx Transaction) IsCoinbase() bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].TxID) == 0 && tx.Inputs[0].Out == -1
}

func (tx Transaction) TrimmedCopy() Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	for _, in := range tx.Inputs {
		inputs = append(inputs, TXInput{
			TxID:      in.TxID,
			Out:       in.Out,
			Signature: nil,
			PublicKey: nil,
		})
	}

	for _, out := range tx.Outputs {
		outputs = append(outputs, out)
	}

	return Transaction{
		ID:      tx.ID,
		Inputs:  inputs,
		Outputs: outputs,
	}
}

func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]Transaction) error {
	if tx.IsCoinbase() {
		return nil
	}

	for _, in := range tx.Inputs {
		if prevTXs[hex.EncodeToString(in.TxID)].ID == nil {
			return errors.New("Previous transactions are invalid")
		}
	}

	txCopy := tx.TrimmedCopy()

	for idx, in := range txCopy.Inputs {
		prevTX := prevTXs[hex.EncodeToString(in.TxID)]
		txCopy.Inputs[idx].Signature = nil // just incase.....
		txCopy.Inputs[idx].PublicKey = prevTX.Outputs[in.Out].PubkeyHash

		hash, err := txCopy.Hash()
		if err != nil {
			return err
		}

		txCopy.ID = hash
		txCopy.Inputs[idx].PublicKey = nil

		r, s, err := ecdsa.Sign(rand.Reader, &privKey, txCopy.ID)
		if err != nil {
			return err
		}

		tx.Inputs[idx].Signature = append(r.Bytes(), s.Bytes()...)
	}

	return nil
}

func (tx *Transaction) Verify(prevTXs map[string]Transaction) (bool, error) {
	if tx.IsCoinbase() {
		return true, nil
	}

	for _, in := range tx.Inputs {
		if prevTXs[hex.EncodeToString(in.TxID)].ID == nil {
			return false, errors.New("Previous transactions are invalid")
		}
	}

	txCopy := tx.TrimmedCopy()
	curve := elliptic.P256()

	for idx, in := range tx.Inputs {
		prevTX := prevTXs[hex.EncodeToString(in.TxID)]

		txCopy.Inputs[idx].Signature = nil // just incase.....
		txCopy.Inputs[idx].PublicKey = prevTX.Outputs[in.Out].PubkeyHash

		hash, err := txCopy.Hash()
		if err != nil {
			return false, err
		}
		txCopy.ID = hash

		r := big.Int{}
		s := big.Int{}
		sigLen := len(in.Signature)
		r.SetBytes(in.Signature[:(sigLen / 2)])
		s.SetBytes(in.Signature[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(in.PublicKey)
		x.SetBytes(in.PublicKey[:(keyLen / 2)])
		y.SetBytes(in.PublicKey[(keyLen / 2):])

		rawPubKey := ecdsa.PublicKey{
			Curve: curve,
			X:     &x,
			Y:     &y,
		}

		if !ecdsa.Verify(&rawPubKey, txCopy.ID, &r, &s) {
			return false, nil
		}

		txCopy.Inputs[idx].PublicKey = nil
	}

	return true, nil
}

func NewUTXOTransaction(receiver string, amount uint, wallet *Wallet, UTXOSet *UTXOSet) (Transaction, error) {
	pubKeyHash, err := HashPubKey(wallet.PubKey)
	if err != nil {
		return Transaction{}, err
	}

	acc, spendableOuts, err := UTXOSet.FindSpendableOutputs(pubKeyHash, amount)
	if err != nil {
		return Transaction{}, err
	}

	if acc < amount {
		return Transaction{}, errors.New("Not enough funds")
	}

	var inputs []TXInput
	var outputs []TXOutput

	// build a list of inputs
	for ID, outs := range spendableOuts {
		txID, err := hex.DecodeString(ID)
		if err != nil {
			return Transaction{}, err
		}

		for _, out := range outs {
			in := TXInput{
				TxID:      txID,
				Out:       out,
				Signature: nil,
				PublicKey: wallet.PubKey,
			}

			inputs = append(inputs, in)
		}
	}

	// build a list of outputs
	outputs = append(outputs, NewTXOutput(amount, receiver))

	sender, err := wallet.GetAddress()
	if err != nil {
		return Transaction{}, err
	}

	if acc > amount {
		// a change
		outputs = append(outputs, NewTXOutput(acc-amount, sender))
	}

	tx := Transaction{
		ID:      nil,
		Inputs:  inputs,
		Outputs: outputs,
	}

	tx.ID, err = tx.Hash()
	if err != nil {
		return Transaction{}, err
	}

	UTXOSet.Blockchain.SignTX(tx, wallet.PrivKey)
	return tx, nil
}
