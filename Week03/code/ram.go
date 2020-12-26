package main

import (
	"fmt"
	"time"
)

type Stu interface {
	Say()
}

type Ben struct {
	Name string
}

func (s *Ben) Say() {
	fmt.Println("Ben say:I am ", s.Name)
}

type Mike struct {
	Name string
}

func (s *Mike) Say() {
	fmt.Println("Mike say: I am ", s.Name)
}

type ZhangSan struct {
	Age  int
	Name string
}

func (s *ZhangSan) Say() {
	fmt.Println("Mike say: I am ", s.Name)
}

//
func foo1() {
	var stu Stu

	b := &Ben{"ben"}
	m := &Mike{"mike"}
	go func() {
		for {
			stu = b
			stu.Say()
		}
	}()

	go func() {
		for {
			stu = m
			stu.Say()
		}
	}()

	time.Sleep(time.Second * 5)
}

// 当内存布局不同时，且发生data race 会panic
// interface 内部是两个指针，也就是16字节，含有两个机器字，赋值操作是非原子的
// type 指向ben，value指向zhangsan，这个时候程序按照ben的内存布局反射value的值，即发生panic
func foo2() {
	var stu Stu

	b := &Ben{"ben"}
	m := &ZhangSan{Name: "zhangsan"}
	go func() {
		for {
			stu = b
			stu.Say()
		}
	}()

	go func() {
		for {
			stu = m
			stu.Say()
		}
	}()

	time.Sleep(time.Second * 5)
}

func main() {
	//foo1()
	foo2()
}
