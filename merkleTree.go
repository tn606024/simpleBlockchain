package simpleBlockchain

import (
	"crypto/sha256"
)

type merkleTree struct {
	root *merkleNode
}

type merkleNode struct {
	left *merkleNode
	right *merkleNode
	hash []byte
}

func CalculateMerkleRoot(transactions []*Transaction) []byte {
	var dataset [][]byte
	for _, tx := range transactions {
		dataset = append(dataset, tx.newHash())
	}
	mt := NewMerkleTree(dataset)
	return mt.root.hash
}

func NewMerkleTree(dataset [][]byte) *merkleTree{
	var merkleNodes []*merkleNode
	if len(dataset) == 1 {
		return &merkleTree{root:NewMerkleNode(nil, nil, dataset[0])}
	}
	if len(dataset) % 2 != 0 {
		dataset = append(dataset, dataset[len(dataset)-1])
	}
	for i:=0 ; i<len(dataset); i=i+1 {
		merkleNodes = append(merkleNodes, NewMerkleNode(nil, nil, dataset[i]))
	}
	for i:=0; i<len(merkleNodes); i=i+2 {
		if i == len(merkleNodes) -1 {
			break
		}
		merkleNodes = append(merkleNodes,NewMerkleNode(merkleNodes[i], merkleNodes[i+1], nil))
	}
	return &merkleTree{root:merkleNodes[len(merkleNodes)-1]}
}

func NewMerkleNode(left, right *merkleNode, data []byte) *merkleNode{
	var node merkleNode
	if left == nil && right == nil {
		hash := sha256.Sum256(data)
		node.hash = hash[:]
	}else{
		hash := sha256.Sum256(append(left.hash, right.hash...))
		node.hash = hash[:]
	}
	node.left = left
	node.right = right
	return &node
}