package main

import (
	"encoding/json"
	"fmt"
)

type User struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Active  bool    `json:"active"`
	Address Address `json:"address"`
}

type Address struct {
	City       string `json:"city"`
	PostalСode string `json:"postal_code"`
}

func main() {
	var User User

	jsondata := `{
		"id": 42,
		"name": "Alex",
		"active": true,
		"address": {
			"city": "Kyoto",
			"postal_code": "600-0000"
		}
	}`

	_ = json.Unmarshal([]byte(jsondata), &User)

	fmt.Println(User)
}
