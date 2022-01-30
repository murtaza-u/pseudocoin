package core

import (
	"bytes"
	"crypto/sha256"
	"log"
	"math"
	"math/big"
	"time"
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
	target = target.Lsh(target, 256-targetBits)
	return &PoW{b, target}
}

func (pow *PoW) prepareData(nonce uint64) []byte {
	var targetBytes, nonceBytes, timeBytes []byte

	targetBytes, err := IntToBytes(targetBits)
	nonceBytes, err = IntToBytes(int64(nonce))
	timeBytes, err = IntToBytes(time.Now().Unix())

	if err != nil {
		log.Panic(err)
	}

	return bytes.Join([][]byte{
		targetBytes,
		nonceBytes,
		timeBytes,
		// pow.block.prevBlockHash,
		// pow.block.hashTXs(),
	}, []byte{})
}

func (pow *PoW) Run() ([]byte, uint64) {
	var hash [32]byte
	var hashInt big.Int
	var nonce uint64

	for nonce < maxNonce {
		hash = sha256.Sum256(pow.prepareData(nonce))
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.Target) == -1 {
			break
		}

		nonce++
	}

	return hash[:], nonce
}

func (pow *PoW) validate() bool {
	var hashInt big.Int
	hash := sha256.Sum256(pow.prepareData(pow.Block.Nonce))
	hashInt.SetBytes(hash[:])
	return hashInt.Cmp(pow.Target) == -1
}
