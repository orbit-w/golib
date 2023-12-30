package utils

import (
	"log"
	"runtime/debug"
)

func GoSafe(handle func()) {
	go func() {
		if r := recover(); r != nil {
			log.Printf("recoverd: %v", r)
			debug.PrintStack()
		}
		handle()
	}()
}
