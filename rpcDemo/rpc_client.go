package main

import (
	"fmt"
	"log"
	"net/rpc"
)

type ArithRequest struct {
	A int
	B int
}

type ArithResponse struct {
	Pro int
	Quo int
	Rem int
}

func main() {
	conn, err := rpc.DialHTTP("tcp", "127.0.0.1:8095")
	if err != nil {
		log.Fatalln(err)
	}
	req := ArithRequest{9, 2}
	var res ArithResponse
	err = conn.Call("Arith.Multiply", req, &res)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%d * %d = %d\n", req.A, req.B, res.Pro)
}
