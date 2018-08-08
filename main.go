package main

import (
	"xu.com/mykv/mykv"
	"fmt"
)
func main() {
	db,_ := mykv.Open("test.txt",0666,nil)
	fmt.Printf("%T",db)
}