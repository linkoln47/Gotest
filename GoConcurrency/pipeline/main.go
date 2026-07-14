package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
)

type Input struct {
	A string
	B string
}

type Cout struct {
	freqCount map[rune]int
}

func GatherAndProcess(ctx context.Context, data Input) (Cout, error) {
	ctx, cancel := context.WithTimeout(ctx, 50*time.Millisecond)
	defer cancel()

	ab := newABProcessor()
	ab.start(ctx, data)
	Cinput, err := ab.wait(ctx)
	if err != nil {
		return Cout{}, err
	}
	c := newCProcessor()
	c.start(ctx, Cinput)
	out, err := c.wait(ctx)
	return out, err
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("expexted input to be processed")
		os.Exit(1)
	}
	cout, err := GatherAndProcess(context.Background(), Input{
		A: os.Args[1],
		B: os.Args[2],
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(cout)
}
