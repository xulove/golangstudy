package main

import (
	"fmt"
	"unsafe"
)

type Person struct {
	name string
	age int
	address string
}

func main() {
	p := &Person{"xu",18,"zhongguo123"}
	//fmt.Println(unsafe.Offsetof(Person{}.address))
	//temp := (*[unsafe.Offsetof(Person{}.address)]byte)(unsafe.Pointer(p))[:]
	//fmt.Println(temp)
	fmt.Println(unsafe.Sizeof(*p))
}