package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync/atomic"
	"time"
)

type Window struct {
	caps        int   // 窗口大小
	currentTime int64 // 当前时间
	ring        *Ring // 环形链表
}

func NewWindow(caps int) *Window {
	return &Window{
		caps: caps,
		ring: New(caps),
	}
}

func (w *Window) Inc() {
	nowTime := time.Now().Unix()
	if nowTime != w.currentTime {
		w.currentTime = nowTime
		w.ring = w.ring.Next()
		w.ring.Value = 0
	}

	atomic.AddInt64(&w.ring.Value, 1)
}

func (w *Window) Len(caps ...int) int64 {
	realCaps := w.caps
	if len(caps) > 0 {
		if caps[0] < realCaps {
			realCaps = caps[0]
		}
	}

	count := int64(0)
	node := w.ring
	for i := 0; i < realCaps; i++ {
		count += node.Value
		node = node.Prev()
	}

	return count
}

// 粘贴的 ring包 将value的类型由interface 改为int
type Ring struct {
	next, prev *Ring
	Value      int64 // for use by client; untouched by this library
}

func (r *Ring) init() *Ring {
	r.next = r
	r.prev = r
	return r
}

// Next returns the next ring element. r must not be empty.
func (r *Ring) Next() *Ring {
	if r.next == nil {
		return r.init()
	}
	return r.next
}

// Prev returns the previous ring element. r must not be empty.
func (r *Ring) Prev() *Ring {
	if r.next == nil {
		return r.init()
	}
	return r.prev
}

// New creates a ring of n elements.
func New(n int) *Ring {
	if n <= 0 {
		return nil
	}
	r := new(Ring)
	p := r
	for i := 1; i < n; i++ {
		p.next = &Ring{prev: p}
		p = p.next
	}
	p.next = r
	r.prev = p
	return r
}

func main() {
	caps := 10
	w := NewWindow(caps)

	rand.Seed(time.Now().UnixNano())
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

label:
	for {
		select {
		case <-ctx.Done():
			break label
		default:
			w.Inc()
			time.Sleep(time.Duration(rand.Intn(100)) * time.Microsecond)
		}
	}

	for i := 1; i < caps+1; i++ {
		fmt.Printf("%d秒 窗口计数为 %d\n", i, w.Len(i))
	}
}
