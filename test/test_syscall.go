package main

import (
	"os"
	"log"
	"syscall"
	"fmt"
)

func main() {
	f,err := os.OpenFile("mmap.bin",os.O_RDWR|os.O_CREATE,0644)
	if err != nil {
		//log.Fatalln(err)
	}

	if _,err := f.WriteAt([]byte{byte(0)},1<<8);nil != err{
		//log.Fatalln(err)
	}
	data,err := syscall.Mmap(int(f.Fd()),0,1<<8,syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		//log.Fatalln(err)
	}
	fmt.Println(string(data))

	if err := f.Close(); nil != err {
		log.Fatalln(err)
	}
	for i, v := range []byte("hello world") {
		data[i] = v
	}
	fmt.Println("here1-",string(data))
	if err := syscall.Munmap(data); nil != err {
		log.Fatalln(err)
	}
}
