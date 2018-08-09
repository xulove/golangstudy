package main

import (
	"fmt"
	"unsafe"
)

type Person struct {
	name string
	age string
}

func main() {
	// 证明数组的指针就是数据中第一个元素的指针
	//arr := [2]string{"hello","world"}
	//fmt.Println(unsafe.Pointer(&arr))
	//fmt.Println(&arr[0])
	//fmt.Println(unsafe.Pointer(&arr[0]))
	//fmt.Println(&arr[1])
	p := new(Person)
	n := (*[10]byte)(unsafe.Pointer(&(p.age)))
	fmt.Println(n)
	fmt.Printf("%T\n",n)
	m := (*[10]byte)(unsafe.Pointer(&(p.age)))[2:]
	fmt.Println(m)
	fmt.Printf("%T",m)

}
