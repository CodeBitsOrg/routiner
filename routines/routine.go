package routines

import (
	"fmt"
	"log"
	"runtime/debug"
	"sync"
)

type Routine struct {
	Input  chan interface{}
	Output chan string
	Wg     *sync.WaitGroup
	Quit   chan int
}

func (r *Routine) Info(str string) {
	r.Output <- str
}

func (r *Routine) Work(obj interface{}) {
	r.Wg.Add(1)

	r.Input <- obj
}

func (r *Routine) Finish() {
	r.Wg.Wait()
	r.Quit <- 0
}

func (r *Routine) Recover() func() {
	f := func() {
		if err := recover(); err != nil {
			log.Print(fmt.Errorf("%s\n%s", err, debug.Stack()))
		}
	}

	return f
}

func runManager(r *Routine, manager func(r *Routine)) {
	// Recover the application in case of a panic.
	defer r.Recover()()

	// Call manager method from the routine's handler.
	manager(r)

	r.Finish()
}

func runWorker(r *Routine, worker func(r *Routine, input interface{})) {
	// Recover the application in case of a panic.
	defer r.Recover()()

	for {
		select {
		case input := <-r.Input:
			worker(r, input)

			r.Wg.Done()
		}
	}
}
