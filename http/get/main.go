package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Todo struct {
	UserID    int    `json:"userId"`
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func Fetch(url string) (Todo, error) {
	var todo Todo
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return Todo{}, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return Todo{}, fmt.Errorf("%s", resp.Status)
	}

	// body, _ := io.ReadAll(resp.Body)

	if err := json.NewDecoder(resp.Body).Decode(&todo); err != nil {
		return Todo{}, err
	}

	return todo, nil
}

func main() {
	url := "https://jsonplaceholder.typicode.com/todos/1"
	todo, err := Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("todo: %+v\n", todo)
}
