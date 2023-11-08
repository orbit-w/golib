package astar

import (
	"github.com/orbit-w/golib/v1/bases/container/heap_list"
	"github.com/orbit-w/golib/v1/bases/misc/number_utils"
)

/*
   @Author: orbit-w
   @File: A_star
   @2023 11月 周一 22:11
*/

type Rec struct {
	X, Y int32
	Ori  *Rec
}

type FindingPath struct {
	open  *heap_list.HeapList[int32, Rec, int32]
	close map[int32]struct{}
	check func(int322 int32)
}

func (f *FindingPath) Search() {

}

func heuristic(sX, sY, eX, eY int32) int32 {
	d := HVMoveCost
	return d * (number_utils.ABS[int32](eX-sX) + number_utils.ABS[int32](eY-sY))
}
