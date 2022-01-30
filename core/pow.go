package core

import (
	"math/big"
)

type PoW struct {
	Block  *Block
	Target *big.Int
}

const targetBits = 16

func NewPoW(b *Block) *PoW {
	target := big.NewInt(1)
	target = target.Lsh(target, 256-targetBits)
	return &PoW{b, target}
}
