package main

import (
	"fmt"
	"time"
)

// recover panic
func Async(foo func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("recover:", err)
			}
		}()

		foo()
	}()
}

func foo(a int, b string) {
	fmt.Println(a)
	time.Sleep(time.Second)
	panic(b)
}

func main() {
	for i := 0; i < 3; i++ {
		Async(func() {
			foo(1, "haha")
		})
	}

	fmt.Println("==============================")
	time.Sleep(time.Second * 2)
}
