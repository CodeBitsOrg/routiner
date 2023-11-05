package routines

import (
	"fmt"
	"sync"
)

func Bootstrap(
	manager func(r *Routine),
	worker func(r *Routine, o interface{}),
) {
	wg := new(sync.WaitGroup)

	input := make(chan interface{})
	output := make(chan string)
	quit := make(chan int)

	r := Routine{
		input,
		output,
		wg,
		quit,
	}

	workers := 2

	for i := 1; i <= workers; i++ {
		go runWorker(&r, worker)
	}

	go runManager(&r, manager)

	for {
		select {
		case message := <-r.Output:
			fmt.Println(message)
		case <-r.Quit:
			fmt.Println("The task has been finished :)")
			close(r.Output)
			return
		}
	}
}
