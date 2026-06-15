package main

import (
	_ "embed"
	"fmt"
	"os"
	"strings"
	"time"
)

//go:embed passwords.txt
var passwords string

func main() {
	start := time.Now()
	defer func() {
		fmt.Println("time: ", time.Since(start))	
	}()
	pwds := strings.Split(passwords, "\n")
	if len(os.Args) > 1 {
		for _, v := range pwds {
			if v == os.Args[1] {
				fmt.Println("true")
				return
			}
		}
		fmt.Println("false")
	}	
}
