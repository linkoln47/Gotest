package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func SendMessage(url string, msg string) (int, string, error) {
	body := strings.NewReader(msg)

	resp, err := http.Post(url, "text/plain", body)
	if err != nil {
		return 0, "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, "", err
	}

	if resp.StatusCode != http.StatusOK {
		return resp.StatusCode, string(respBody), fmt.Errorf("%s", resp.Status)
	}

	return resp.StatusCode, string(respBody), nil
}

func main() {
	url := "https://httpbin.org/post"
	status, body, err := SendMessage(url, "Hello World")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(status)
	fmt.Println(body)
}
