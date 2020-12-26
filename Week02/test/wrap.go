package main

import (
	"errors"
	"fmt"

	xerror "github.com/pkg/errors"
)

func foo1() error {
	err := db()
	return xerror.Wrap(err, "foo1")
	//return xerror.WithMessage(err, "foo1")
}

func foo2() error {
	err := foo1()
	return xerror.WithMessage(err, "foo2")
	//return xerror.WithMessage(err, "foo1")
}

func db() error {
	return errors.New("db error")
}

func main() {
	err := foo2()
	fmt.Printf("%+v\n", err)
}
