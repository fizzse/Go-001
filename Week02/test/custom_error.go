package main

import (
	"log"
)

type NotFound interface {
	notFound() bool
}

type MyError struct {
	msg string
}

func NewMyError(msg string) error {
	return &MyError{msg: msg}
}

func (e *MyError) Error() string {
	return e.msg
}

func (e *MyError) notFound() bool {
	return e.msg == "not found"
}

func IsNotFound(err error) bool {
	if myErr, ok := err.(*MyError); ok {
		return myErr.notFound()
	}

	return false
}

func mockNotFoundError() error {
	return NewMyError("not found")
}

func main() {
	err := mockNotFoundError()
	if err != nil {
		if IsNotFound(err) {
			log.Println("not found")
			return
		}

		log.Fatal("some error")
	}

	log.Println("success")
}
