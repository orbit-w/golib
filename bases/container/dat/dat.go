package dat

import (
	"fmt"
	"golib/bases/misc/math"
	"sort"
)

/*
   @Time: 2023/8/22 00:10
   @Author: david
   @File: dat
*/

type DAT struct {
	size  int         //容量
	base  []int       //转移基数
	check []int       //dat 映射父子节点唯一关系性
	trie  *Trie       //trie树结构
	ks    strKeySlice //关键词集合
}

func (ins *DAT) Build(keywords strKeySlice) error {
	ins.ks = keywords
	sort.Sort(ins.ks)

	lt := new(Trie)
	lt.Fetch(ins)
	ins.trie = lt
	ins.base = make([]int, 1024)
	ins.check = make([]int, 1024)
	ins.base[0] = RootState
	ins.PriorityInsert()
	return nil
}

func (ins *DAT) Find(keyword []rune) bool {
	var index int
	for _, r := range keyword {
		pos := ins.pos(ins.state(index), r)
		if ins.check[pos] != index {
			return false
		}
		index = pos
	}
	return ins.base[index] < 0
}

func (ins *DAT) PriorityInsert() {
	ins.priorityInsert(ins.trie.root)
	for i := range ins.trie.linkedList {
		siblings := ins.trie.linkedList[i].siblings
		for j := range siblings {
			ins.priorityInsert(siblings[j])
		}
	}
}

func (ins *DAT) priorityInsert(father *Node) {
	if len(father.children) == 0 {
		return
	}
	//确定最佳转移基数
	k := ins.state(father.index)
	for {
	COMPLETE:
		for i := range father.children {
			node := father.children[i]
			pos := k + int(node.code)
			if ins.base[pos] != 0 {
				k++
				goto COMPLETE
			}
		}
		break
	}

	ins.setState(father.index, k)
	for i := range father.children {
		node := father.children[i]
		pos := k + int(node.code)
		//记录节点在base 中 index
		node.index = pos
		//记录父子节点关系
		ins.check[pos] = father.index
		if node.isLeaf {
			ins.base[pos] = -k
		} else {
			ins.base[pos] = k
		}
	}
}

//TODO: 调整容量方式？如何小范围调整不合适？
func (ins *DAT) resize(size int) {
	if size <= ins.size {
		return
	}

	//deep copy
	newBase := make([]int, size)
	copy(newBase, ins.base)
	ins.base = newBase

	newCheck := make([]int, size)
	copy(newCheck, ins.check)
	ins.check = newCheck

	ins.size += size
}

func (ins *DAT) pos(state int, b rune) int {
	return state + int(b)
}

func (ins *DAT) state(index int) int {
	return math.ABS[int](ins.base[index])
}

func (ins *DAT) setState(index, state int) {
	remain := ins.base[index]
	if remain < 0 {
		ins.base[index] = -state
	} else {
		ins.base[index] = state
	}
}

func Print(dat *DAT) {
	fmt.Println("base: ", dat.base)
	fmt.Println("check: ", dat.check)
}
