package core

import "crypto/sha256"

type MerkleTree struct {
	Root *MerkleNode
}

type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Data  []byte
}

func newMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode {
	node := MerkleNode{}
	var hash [32]byte

	if left == nil && right == nil {
		hash = sha256.Sum256(data)
	} else {
		prevHashes := append(left.Data, right.Data...)
		hash = sha256.Sum256(prevHashes)
	}

	node.Left = left
	node.Right = right
	node.Data = hash[:]

	return &node
}

func NewMerkleTree(data [][]byte) *MerkleTree {
	var nodes []MerkleNode

	if len(data)%2 != 0 {
		data = append(data, data[len(data)-1])
	}

	for _, datum := range data {
		node := newMerkleNode(nil, nil, datum)
		nodes = append(nodes, *node)
	}

	for i := 0; i < len(data)/2; i++ {
		var newLevel []MerkleNode

		for j := 0; j < len(nodes); j += 2 {
			node := newMerkleNode(&nodes[j], &nodes[j+1], nil)
			newLevel = append(newLevel, *node)
		}

		nodes = newLevel
	}

	tree := MerkleTree{&nodes[0]}

	return &tree
}
