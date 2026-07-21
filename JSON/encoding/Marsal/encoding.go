package main

import (
	"encoding/json"
	"fmt"
)

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	InStock     bool    `json:"in_stock"`
	Description string  `json:"description,omitempty"`
}

func main() {
	product := Product{
		ID:      10,
		Name:    "Keyboard",
		Price:   5990.50,
		InStock: true,
	}

	data, _ := json.Marshal(product)

	fmt.Printf("%s\n", string(data))
}
