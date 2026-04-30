package main

import (
	"errors"
	"fmt"
)

type Person struct {
	FirstName string
	LastName  string
	Age       int
}

func ValidatePerson(p Person) error {
	var err1, err2, err3 error
	if len(p.FirstName) == 0 {
		err1 = errors.New("field FirstName cannot be empty")
	}
	if len(p.LastName) == 0 {
		err2 = errors.New("field LastName cannot be empty")
	}
	if p.Age < 0 {
		err3 = errors.New("field Age cannot be negative")
	}
	if err1 != nil || err2 != nil || err3 != nil {
		err := fmt.Errorf("validation failed with errors: %w, %w, %w", err1, err2, err3)
		return err
	}
	return nil
}

func main() {
	p := Person{
		FirstName: "",
		LastName:  "",
		Age:       -1,
	}
	err := ValidatePerson(p)
	if err != nil {
		fmt.Printf("validation failed with error: %v\n", err)
	}
}
