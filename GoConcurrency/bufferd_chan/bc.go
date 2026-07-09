package main

import "fmt"

const conc = 10

func processChannel(ch chan int) []int {
	results := make(chan int, conc)
	for i := 0; i < conc; i++ {
		go func() {
			v := <-ch
			results <- process(v)
		}()
	}
	var out []int
	for i := 0; i < conc; i++ {
		out = append(out, <-results)
	}
	return out
}

func process(i int) int {
	return i * 2
}

func main() {
	vals := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	ch := make(chan int, conc)
	go func() {
		for _, v := range vals {
			ch <- v
		}
	}()
	result := processChannel(ch)
	fmt.Println(result)
}
