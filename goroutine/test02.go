package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go watch(ctx, "watch_1")
	go watch(ctx, "watch_2")
	go watch(ctx, "watch_3")
	time.Sleep(10 * time.Second)
	fmt.Println("now ,you can stop all!")
	cancel()
	time.Sleep(5 * time.Second)
	fmt.Println("over")
}

func watch(ctx context.Context, str string) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println(str, " quit")
			return
		default:
			fmt.Println(str)
			time.Sleep(2 * time.Second)
		}
	}
}
