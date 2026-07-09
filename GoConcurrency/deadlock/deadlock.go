package main

import "fmt"

func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)
	done := make(chan struct{})

	go func() {
		inRoutine := 1
		ch1 <- inRoutine
		fromMain := <-ch2
		fmt.Println("goroutine: ", inRoutine, fromMain)
		close(done)
	}()
	inMain := 2
	fromRoutine := <- ch1
	ch2 <- inMain
	fmt.Println("main: ", fromRoutine, inMain)
	<- done
}
