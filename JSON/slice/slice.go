package main

import (
	"encoding/json"
	"fmt"
)

type Language struct {
	Name  string `json:"name"`
	Score int    `json:"score"`
}

func main() {
	var languages []Language
	data := `[
	{"name":"Go","score":90},
	{"name":"Python","score":75},
	{"name":"Rust","score":85},
	{"name":"Java","score":60}
	]`

	_ = json.Unmarshal([]byte(data), &languages)

	fmt.Println(languages)

	max, sum := 0, 0
	var maxLanguage string
	filtered := make([]Language, 0)

	for _, v := range languages {
		sum += v.Score
		if max < v.Score {
			max = v.Score
			maxLanguage = v.Name
		}
		if v.Score >= 80 {
			filtered = append(filtered, v)
		}
	}
	fmt.Println(sum / len(languages))
	fmt.Println(maxLanguage)

	newdata, _ := json.MarshalIndent(filtered, "", " ")

	fmt.Println(string(newdata))
}
