package main

import (
	"fmt"
	"sync"
)

func count() {
	wg := sync.WaitGroup{}
	wg.Add(2)

	num := 0
	go func() {
		defer wg.Done()
		for i := 0; i < 10000; i++ {
			num++
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 10000; i++ {
			num++
		}
	}()

	wg.Wait()

	fmt.Println("num:", num)
}

func countWithLock() {
	lock := sync.Mutex{}

	wg := sync.WaitGroup{}
	wg.Add(2)

	num := 0
	go func() {
		defer wg.Done()
		for i := 0; i < 100000; i++ {
			lock.Lock()
			num++
			lock.Unlock()
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 100000; i++ {
			lock.Lock()
			num++
			lock.Unlock()
		}
	}()

	wg.Wait()

	fmt.Println("count with lock num:", num)
}

func main() {
	count()
	countWithLock()
}
