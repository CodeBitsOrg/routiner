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
		go r.runWorker(worker)
	}

	go r.runManager(manager)

	for {
		select {
		case message := <-r.Output:
			log.Println(message)
		case <-r.Quit:
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

// Recover the application in case of a panic.
func (r *Routiner) Recover() func() {
	f := func() {
		if err := recover(); err != nil {
			log.Print(fmt.Errorf("%s\n%s", err, debug.Stack()))
		}
	}

	return f
}

func (r *Routiner) runManager(manager func(r *Routiner)) {
	defer r.Recover()()

	manager(r)

	r.Finish()
}

func (r *Routiner) runWorker(worker func(r *Routiner, input interface{})) {
	defer r.Recover()()

	for {
		select {
		case input := <-r.Input:
			worker(r, input)

			r.Wg.Done()
		}
	}
}
