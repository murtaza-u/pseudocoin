package core

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

type PoW struct {
	Block  *Block
	Target *big.Int
}

const (
	targetBits = 16
	maxNonce   = math.MaxUint64
)

func NewPoW(b *Block) *PoW {
	target := big.NewInt(1)
	target = target.Lsh(target, uint(256-targetBits))
	return &PoW{b, target}
}

func (pow *PoW) PrepareData(nonce uint64) ([]byte, error) {
	var targetBytes, nonceBytes, timeBytes []byte

	targetBytes, err := IntToBytes(int64(targetBits))
	if err != nil {
		return nil, err
	}

	nonceBytes, err = IntToBytes(int64(nonce))
	if err != nil {
		return nil, err
	}

	timeBytes, err = IntToBytes(pow.Block.Timestamp)
	if err != nil {
		return nil, err
	}

	// hashedTXs, err := pow.Block.HashTXs()
	// if err != nil {
	// 	return nil, err
	// }

	return bytes.Join([][]byte{
		pow.Block.PrevBlockHash,
		// hashedTXs,
		timeBytes,
		targetBytes,
		nonceBytes,
	}, []byte{}), err
}

func (pow *PoW) Run() ([]byte, uint64, error) {
	var err error
	var hash [32]byte
	var hashInt big.Int
	var nonce uint64

	for nonce < maxNonce {
		var data []byte
		data, err = pow.PrepareData(nonce)
		if err != nil {
			break
		}

		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)

		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.Target) == -1 {
			break
		}

		nonce++
	}

	fmt.Println()

	return hash[:], nonce, err
}

func (pow *PoW) Validate() (bool, error) {
	var hashInt big.Int
	data, err := pow.PrepareData(pow.Block.Nonce)
	if err != nil {
		return false, err
	}

	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])
	return hashInt.Cmp(pow.Target) == -1, nil
}
