package main

import (
	"fmt"
	"time"

	"github.com/CodeBitsOrg/routiner/routiner"
)

func main() {
	r := routiner.Routiner{
		Workers: 4,
	}

	manager := func(r *routiner.Routiner) {
		for i := 1; i <= 4; i++ {
			r.Work(inputObject{ID: i})
		}
	}

	worker := func(r *routiner.Routiner, o interface{}) {
		obj := o.(inputObject)
		time.Sleep(time.Second)
		r.Info(fmt.Sprintf("Worker %d", obj.ID))
	}

	r.Run(manager, worker)

	fmt.Println("The task has been finished :)")
}

type inputObject struct {
	ID int
}
