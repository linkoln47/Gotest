package main

import "fmt"

type Logger interface {
	Log(message string)
}

type LoggerAdapter func(message string)

func (lg LoggerAdapter) Log(message string) {
	lg(message)
}

func LogOutput(message string) {
	fmt.Println(message)
}
