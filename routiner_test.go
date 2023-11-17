package routiner

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	// Redirecting log output to a buffer, so we can test it.
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(nil)
	}()

	numberOfOperations := 4

	r, manager, worker := scaffold()

	r.Run(manager, worker)

	logOutput := buf.String()

	for i := 1; i <= numberOfOperations; i++ {
		if !strings.Contains(logOutput, fmt.Sprintf("Worker %d", i)) {
			t.Errorf("Worker %d was not found in the log output", i)
		}
	}
}

// Initialize Routiner and its manager and worker functions.
func scaffold() (Routiner, func(r *Routiner), func(r *Routiner, o interface{})) {
	r := Routiner{
		Workers: 4,
	}

	manager := func(r *Routiner) {
		for i := 1; i <= 4; i++ {
			r.Work(i)
		}
	}

	worker := func(r *Routiner, o interface{}) {
		id := o.(int)
		r.Info(fmt.Sprintf("Worker %d", id))
	}

	return r, manager, worker
}
