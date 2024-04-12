package utils

import (
	"fmt"
	"runtime/debug"
)

/*
   @Author: orbit-w
   @File: panic_utils
   @2024 4月 周五 17:57
*/

func RecoverPanic(handle func()) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
			fmt.Println("Stack trace:")
			debug.PrintStack()
		}
	}()

	handle()
}
