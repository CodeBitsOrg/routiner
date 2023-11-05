package test

import (
	"fmt"
	"time"

	"github.com/CodeBitsOrg/routiner/routines"
)

type inputObject struct {
	ID int
}

func Run() {
	routines.Bootstrap(manager, worker)
}

func manager(r *routines.Routine) {
	for i := 1; i <= 4; i++ {
		r.Work(inputObject{ID: i})
	}
}

func worker(r *routines.Routine, o interface{}) {
	obj := o.(inputObject)
	time.Sleep(time.Second)
	r.Info(fmt.Sprintf("Worker %d \n*****", obj.ID))
}
