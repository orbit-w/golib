package astar

/*
   @Author: orbit-w
   @File: base
   @2023 11月 周二 17:45
*/

const (
	HVMoveCost = int32(10) //水平移动代价
	DMoveCost  = int32(14) //斜向移动代价
)

type direction struct {
	h, v int32
}

var (
	zero = struct{}{}

	directions = []direction{
		{1, 0},
		{1, -1},
		{0, -1},
		{-1, -1},
		{-1, 0},
		{-1, 1},
		{0, 1},
		{1, 1},
	}
)
