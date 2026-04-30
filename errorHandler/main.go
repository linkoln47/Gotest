package main

import (
	"errors"
	"fmt"
	"os"
)

type Status int

const (
	InvalidLogin Status = iota + 1
	NotFound
)

type StatusErr struct {
	Status  Status
	Message string
}

func (se StatusErr) Error() string {
	return se.Message
}

func LoginAndGetData(uid, pwd, file string) ([]byte, error) {
	token, err := login(uid, pwd)
	if err != nil {
		return nil, StatusErr{
			Status:  InvalidLogin,
			Message: fmt.Sprintf("invalid login for user %s", uid),
		}
	}
	data, err := getData(token, file)
	if err != nil {
		return nil, StatusErr{
			Status:  NotFound,
			Message: fmt.Sprintf("file %s not found", file),
		}
	}
	return data, nil
}

func login(uid, pwd string) (string, error) {
	if uid == "admin" && pwd == "admin" {
		return "user:admin", nil
	}
	return "", errors.New("bad User")
}

func getData(token, file string) ([]byte, error) {
	if token == "user:admin" {
		switch file {
		case "secret.txt":
			return []byte("pswd aplenty!"), nil
		case "payroll.csv":
			return []byte("everyone's salary"), nil
		}
	}

	return nil, os.ErrNotExist
}

func GenErrBroken(flag bool) error {
	var genErr error
	if flag {
		genErr = StatusErr{
			Status: NotFound,
		}
	}
	return genErr
}

func main() {
	err := GenErrBroken(true)
	fmt.Println("GenErrBroken(true) returns non-nill err:", err != nil)
	err = GenErrBroken(false)
	fmt.Println("GenErrBroken(false) returns non-nill err:", err != nil)
}
