package main

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type Language struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	var b bytes.Buffer
	enc := json.NewEncoder(&b)

	data := []Language{
		{Name: "Fred", Age: 40},
		{Name: "Fred", Age: 40},
		{Name: "Fred", Age: 40},
	}

	for _, language := range data {
		if err := enc.Encode(language); err != nil {
			panic(err)
		}
	}

	fmt.Println(b.String())
}
