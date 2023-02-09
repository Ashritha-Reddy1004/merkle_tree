package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

var leaves = []*Leaf{}
var merkleroot *Node

type Leaf struct {
	Input []byte
	Hash  []byte
}

type Node struct {
	Left  *Node
	Right *Node
	Hash  []byte
}

func AddLeaf(input []byte) *Leaf {
	hash32 := sha256.Sum256(input)
	hash := []byte{}
	hash = append(hash, hash32[:]...)
	return &Leaf{
		Input: input,
		Hash:  hash[:],
	}
}
func AddParent(leaf *Leaf) *Node {
	return &Node{
		Hash: leaf.Hash,
	}
}
func AddNode(left, right *Node) *Node {
	var hash []byte
	if right == nil {
		hash32 := sha256.Sum256(append(left.Hash, left.Hash...))
		hash = append(hash, hash32[:]...)
	} else {
		hash32 := sha256.Sum256(append(left.Hash, right.Hash...))
		hash = append(hash, hash32[:]...)
	}
	return &Node{
		Left:  left,
		Right: right,
		Hash:  hash[:],
	}
}
func Input(input string) {
	leaf := AddLeaf([]byte(input))
	leaves = append(leaves, leaf)
	fmt.Println(input, " is added")
	MerkleTree()
}
func DeleteNode(input string) {
	bytedata := []byte(input)
	for i := 0; i < len(leaves); i++ {
		leafdata := leaves[i].Input
		if bytes.Equal(leafdata, bytedata) {
			leaves = append(leaves[:i], leaves[i+1:]...)
			MerkleTree()
			fmt.Println(input, "is deleted")
			break
		}
		fmt.Println("Not Found")
	}
}
func Check(input string) bool {
	var hash []byte
	bytedata := []byte(input)
	hash32 := sha256.Sum256(bytedata)
	hash = append(hash, hash32[:]...)
	return CheckNode(merkleroot, hash)
}
func CheckNode(root *Node, target []byte) bool {
	if root == nil {
		return false
	}
	if bytes.Equal(root.Hash, target) {
		return true
	}
	var left, right bool
	if root.Left != nil {
		left = CheckNode(root.Left, target)
	}
	if root.Right != nil {
		right = CheckNode(root.Right, target)
	}
	return left || right
}

func MerkleTree() {
	var nodes []*Node
	for i := 0; i < len(leaves); i += 2 {
		var leftnode, rightnode *Node
		leftnode = AddParent(leaves[i])
		if i+1 < len(leaves) {
			rightnode = AddParent(leaves[i+1])
		}
		nodes = append(nodes, AddNode(leftnode, rightnode))
	}

	for len(nodes) > 1 {
		var tempnodes []*Node
		for i := 0; i < len(nodes); i += 2 {
			var leftnode, rightnode *Node
			leftnode = nodes[i]
			if i+1 < len(nodes) {
				rightnode = nodes[i+1]
			}
			node := AddNode(leftnode, rightnode)
			tempnodes = append(tempnodes, node)
		}
		nodes = tempnodes
		fmt.Println(len(nodes))
	}

	if len(nodes) > 0 {
		merkleroot = nodes[0]
	} else {
		merkleroot = nil
	}
}

func Leaves() []*Leaf {
	return leaves
}

func MerkleRoot() *Node {
	return merkleroot
}

func main() {
	Input("40")
	Input("50")
	Input("10")
	initialroot := hex.EncodeToString(MerkleRoot().Hash)
	fmt.Println(initialroot)
	Input("60")
	Input("100")
	Input("30")
	Input("70")
	fmt.Println(Leaves())
	fmt.Println(Check("100"))
	DeleteNode("100")
	fmt.Println(Leaves())
	fmt.Println(Check("50"))
	DeleteNode("50")
	fmt.Println(Leaves())
	updatedroot := hex.EncodeToString(MerkleRoot().Hash)
	fmt.Println(updatedroot)
}
