package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

type Language struct {
	Name  string `json:"name"`
	Score int    `json:"score"`
}

func main() {
	data := `
	{"name":"Go","score":90}
	{"name":"Python","score":75}
	{"name":"Rust","score":85}
	{"name":"Java","score":60}
	`

	dec := json.NewDecoder(strings.NewReader(data))

	for {
		var languages Language
		err := dec.Decode(&languages)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			panic(err)
		}
		fmt.Println(languages.Name, languages.Score)
	}

}
