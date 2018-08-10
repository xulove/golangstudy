package main

import (
	"runtime"
	"fmt"
)

func init() {
	runtime.GOMAXPROCS(1)
}

func say(s string){
	for i := 0; i < 2; i++ {
		runtime.Gosched()
		fmt.Println(s)
	}
}
func main() {
	go say("world")
	say("hello")
}
