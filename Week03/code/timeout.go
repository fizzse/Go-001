package main

import (
	"context"
	"fmt"
	"time"
)

func doWorkWithTimeOut(ctx context.Context, timeout int, foo func()) error {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		return fmt.Errorf("time out")

	case <-func() chan struct{} {
		done := make(chan struct{})
		go func() {
			defer close(done)
			foo()
		}()
		return done
	}():
		return nil
	}
}

func main() {
	doWorkWithTimeOut(context.Background(), 3, func() {
		time.Sleep(10 * time.Second)
	})

	fmt.Println("done")
}
