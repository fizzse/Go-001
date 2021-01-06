package work

import (
	"container/ring"
	"fmt"
)

type Window struct {
}

func foo() {
	r := ring.New(10)
	r.Do(func(i interface{}) {
		fmt.Print(i)
	})
}
