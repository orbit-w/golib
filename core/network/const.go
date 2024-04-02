package network

/*
   @Author: orbit-w
   @File: const
   @2024 4月 周二 15:33
*/

const (
	TypeWorking = 1
	TypeStopped = 2
)

const (
	MaxIncomingPacket = 1<<18 - 1
	HeadLen           = 4 //包头字节数
)
