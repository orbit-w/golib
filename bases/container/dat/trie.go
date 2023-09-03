package dat

import "fmt"

/*
   @Time: 2023/8/22 00:06
   @Author: david
   @File: trie
*/

type Node struct {
	isLeaf   bool
	code     rune
	depth    int
	index    int
	children []*Node
}

func (ins *Node) findChild(r rune) *Node {
	for _, child := range ins.children {
		if child.code == r {
			return child
		}
	}
	return nil
}

type LinkedSiblingNodes struct {
	siblings []*Node
}

type Trie struct {
	dat        *DAT
	root       *Node
	linkedList []LinkedSiblingNodes
}

func (ins *Trie) Fetch(_dat *DAT) {
	ins.dat = _dat
	ins.linkedList = make([]LinkedSiblingNodes, ins.dat.ks.Max())
	ins.root = new(Node)

	for _, key := range ins.dat.ks {
		node := ins.root
		for i, r := range key {
			child := node.findChild(r)
			linked := &ins.linkedList[i]
			if child == nil {
				child = &Node{
					code:   r,
					depth:  node.depth + 1,
					isLeaf: i == len(key)-1,
				}
				node.children = append(node.children, child)
				linked.siblings = append(linked.siblings, child)
			}
			node = child
		}
	}
}

func (ins *Trie) Print() {
	printTree(ins.root, 0)
}

func printTree(node *Node, lv int) {
	for i := 0; i < lv; i++ {
		fmt.Print("  ")
	}
	fmt.Println(string(node.code))
	for _, child := range node.children {
		printTree(child, lv+1)
	}
}
