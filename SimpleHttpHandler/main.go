package main

import (
	"net/http"
)

func main() {
	l := LoggerAdapter(LogOutput)
	ds := NewSimpleDataStore()
	logic := NewSimpleLogic(l, ds)
	c := NewController(l, logic)
	http.HandleFunc("/hello", c.SayHello)
	http.HandleFunc("/bye", c.SayGoodbye)
	http.ListenAndServe(":8080", nil)
}
