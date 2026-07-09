package main

import (
	"fmt"
)

func main() {
	in1 := make(chan int)
	in2 := make(chan int)
	go func() {
		for i := 10; i < 100; i += 10 {
			in1 <- i
		}
		close(in1)
	}()

	go func() {
		for i := 20; i >= 0; i-- {
			in2 <- i
		}
		close(in2)
	}()
	result := readFromTwoChannels(in1, in2)
	fmt.Println(result)
}

func readFromTwoChannels(in1, in2 chan int) []int {
	var out []int
	for count := 0; count < 2; {
		select {
		case v, ok := <-in1:
			if !ok {
				in1 = nil // this case is blocked
				count++
				continue
			}
			out = append(out, v) // append al res
		case v, ok := <-in2:
			if !ok {
				in2 = nil
				count++
				continue
			}
			out = append(out, v)
		}
	}
	return out
}
