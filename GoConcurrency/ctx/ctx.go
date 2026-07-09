package main

import (
	"context"
	"fmt"
)

func count2(ctx context.Context, max int) <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for i := 0; i < max; i++ {
			select {
			case <-ctx.Done():
				return
			case ch <- i:
			}
			// close(ch) -- goroutine leak, will never reach this place
		}
	}()
	return ch
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ch := count2(ctx, 10)
	for i := range ch {
		if i > 5 {
			cancel()
			break
		}
		fmt.Println(i)
	}
}
