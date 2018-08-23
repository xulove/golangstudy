package main

import "fmt"
import (
	"context"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				time.Sleep(1 * time.Second)
				fmt.Println("monitor stoped")
				return

			default:
				fmt.Println("goroutine started")
				time.Sleep(2 * time.Second)
			}
		}
	}(ctx)

	time.Sleep(2 * time.Second)
	fmt.Println("now,you can stop the goroutine")
	cancel()
	fmt.Println("canceling...")
	time.Sleep(5 * time.Second)
	fmt.Println("over")
}
