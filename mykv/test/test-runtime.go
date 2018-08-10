package main

import (
	"fmt"
	"runtime"
)

func init()  {
	fmt.Println("--init--")
	runtime.GOMAXPROCS(1)
}
func main() {
	//打印当前系统的cpu数
	fmt.Println("numcpu:",runtime.NumCPU())
	//打印当前的goroot
	fmt.Println("Goroot:",runtime.GOROOT())
	//打印当前操作系统
	fmt.Println("goos:",runtime.GOOS)
}
