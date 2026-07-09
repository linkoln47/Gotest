package main 

import (
	"fmt"
	"time"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		count("sheep")
		count("fish")

		wg.Done()
	}()

	wg.Wait()
}

func count(thing string) {
	for i:= 0; i < 5; i++ {
		fmt.Println(i, thing)
		time.Sleep(time.Millisecond * 500)
	}
}
