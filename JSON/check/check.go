package main

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type Config struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

func DecodeConfig(r io.Reader) (Config, error) {
	var config Config

	decoder := json.NewDecoder(r)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&config); err != nil {
		return Config{}, err
	}
	var extra json.RawMessage

	err := decoder.Decode(&extra)

	if err == io.EOF {
		return config, nil
	}

	if err != nil {
		return Config{}, err
	}

	return Config{}, err
}

func main() {
	data := `{"host":"localhost","port":8080}`

	config, err := DecodeConfig(strings.NewReader(data))
	if err != nil {
		fmt.Println("decode error: ", err)
		return
	}

	fmt.Println("Host: ", config.Host)
	fmt.Println("Port: ", config.Port)
}
