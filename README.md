## Basic usage

- First define a Routiner object.

```golang
func main() {
	r := routiner.Routiner{
		Workers: 4,
	}
}
```

- Then define manager & worker clousers that will do the job.

```golang
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
```

- Run the job.

```golang
r.Run(manager, worker)

fmt.Println("The task has been finished :)")
```

- Because the second parameter in the worker clouser is designed to accept some sort of an object, please create it also.

```golang
type inputObject struct {
	ID int
}
```