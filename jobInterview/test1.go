package main
import "fmt"
//1、写出下面代码输出内容。
func main(){
	defer_call()
}

func defer_call(){
	defer func(){ fmt.Println("打印前") }()
	defer func(){ fmt.Println("打印中") }()
	defer func(){ fmt.Println("打印后") }()
	
	panic("触发异常")
}