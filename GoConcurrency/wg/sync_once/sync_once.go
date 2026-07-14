package main

import (
	"fmt"
	"sync"
)

type SlowParser interface {
	Parse(string) string
}

type SCPI struct{}

var getParser = sync.OnceValue(func() SlowParser {
	return initParser()
})

func Parse(data2Parse string) string {
	return getParser().Parse(data2Parse)
}

func initParser() SlowParser {
	fmt.Println("initialiazing...")
	return SCPI{}
}

func (s SCPI) Parse(in string) string {
	if len(in) > 1 {
		return in[0:1]
	}
	return ""
}

func main() {
	res1 := Parse("hello")
	fmt.Println(res1)
	res2 := Parse("goodbye")
	fmt.Println(res2)
}
