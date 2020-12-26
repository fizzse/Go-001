package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"golang.org/x/sync/errgroup"
)

func do(exitChan chan struct{}) error {
	select {
	case <-exitChan:
		return fmt.Errorf("time out")
	default:
		time.Sleep(10 * time.Second)
		log.Println("work end")
	}

	return nil
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	group, _ := errgroup.WithContext(ctx)
	group.Go(func() error {

		select {
		case <-ctx.Done():
			return fmt.Errorf("time out")

		case <-func() chan struct{} {
			done := make(chan struct{})
			go func() {
				defer close(done)
				time.Sleep(time.Second * 10)
				fmt.Println("real do done!!!")
			}()
			return done
		}():

			log.Println("work done")
			return nil
		}
	})

	log.Println("work start")

	if err := group.Wait(); err != nil {
		log.Println("error", err)
	}

	log.Println("end")

	time.Sleep(time.Second * 10)
}
