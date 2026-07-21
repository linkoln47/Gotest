package main

import (
	"encoding/json"
	"fmt"
)

type UpdateUserRequest struct {
	Name   *string `json:"name"`
	Age    *int    `json:"age"`
	Active *bool   `json:"active"`
}

func main() {

	data := `{
		"name": "Mika",
		"active": false
	}`

	var u UpdateUserRequest
	var fields map[string]json.RawMessage

	_ = json.Unmarshal([]byte(data), &u)

	fmt.Println(*u.Name)

	_ = json.Unmarshal([]byte(data), &fields)
	for name, value := range fields {
		fmt.Printf("%s was provided: %s\n", name, value)
	}

}
