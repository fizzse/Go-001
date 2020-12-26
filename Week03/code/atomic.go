package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type Config struct {
	value []int
}

func foo3() {
	c := &Config{}
	go func() {
		i := 0
		for {
			i++
			c.value = []int{i, i + 1, i + 2, i + 3, i + 4, i + 5}
		}
	}()

	wg := sync.WaitGroup{}
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			for i := 0; i < 1000; i++ {
				fmt.Println(c.value)
			}
			wg.Done()
		}()
	}

	wg.Wait()
}

func foo4() {
	c := &Config{}
	v := &atomic.Value{}
	v.Store(c)
	go func() {
		i := 0
		for {
			i++
			cfg := v.Load().(*Config)
			cfg.value = []int{i, i + 1, i + 2, i + 3, i + 4, i + 5}
			v.Store(cfg)
		}
	}()

	wg := sync.WaitGroup{}
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			for i := 0; i < 100; i++ {
				cfg := v.Load().(*Config)
				fmt.Println(cfg.value)
			}
			wg.Done()
		}()
	}

	wg.Wait()
}

func main() {
	//foo3()
	foo4()
}
