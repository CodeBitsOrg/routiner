package routiner

import (
	"fmt"
	"log"
	"runtime/debug"
	"sync"
)

type Routiner struct {
	Input   chan interface{}
	Output  chan string
	Wg      *sync.WaitGroup
	Quit    chan int
	Workers int
}

func (r *Routiner) Run(
	manager func(r *Routiner),
	worker func(r *Routiner, o interface{}),
) {
	wg := new(sync.WaitGroup)

	input := make(chan interface{})
	output := make(chan string)
	quit := make(chan int)

	r.Input = input
	r.Output = output
	r.Wg = wg
	r.Quit = quit

	if r.Workers == 0 {
		r.Workers = 1
	}

	for i := 1; i <= r.Workers; i++ {
		go runWorker(r, worker)
	}

	go runManager(r, manager)

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

func (r *Routiner) Info(str string) {
	r.Output <- str
}

func (r *Routiner) Work(obj interface{}) {
	r.Wg.Add(1)

	r.Input <- obj
}

func (r *Routiner) Finish() {
	r.Wg.Wait()
	r.Quit <- 0
}

func (r *Routiner) Recover() func() {
	f := func() {
		if err := recover(); err != nil {
			log.Print(fmt.Errorf("%s\n%s", err, debug.Stack()))
		}
	}

	return f
}

func runManager(r *Routiner, manager func(r *Routiner)) {
	// Recover the application in case of a panic.
	defer r.Recover()()

	// Call manager method from the routiner's handler.
	manager(r)

	r.Finish()
}

func runWorker(r *Routiner, worker func(r *Routiner, input interface{})) {
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
