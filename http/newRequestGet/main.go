package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type Comment struct {
	PostID int    `json:"postId"`
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Body   string `json:"body"`
}

func Fetch(url string, postID int) ([]Comment, error) {
	var c []Comment

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest(
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	query.Set("postId", strconv.Itoa(postID))
	req.URL.RawQuery = query.Encode()

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Go-http-comments")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s", resp.Status)
	}
	if err := json.NewDecoder(resp.Body).Decode(&c); err != nil {
		return nil, err
	}

	return c, nil
}

func main() {
	url := "https://jsonplaceholder.typicode.com/comments"
	comments, err := Fetch(url, 1)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("count: ", len(comments))

	for _, comm := range comments {
		fmt.Printf("%d: %s\n", comm.ID, comm.Email)
	}
}
