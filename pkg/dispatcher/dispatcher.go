package dispatcher

import (
	"fmt"
	"math/rand"

	"github.com/google/uuid"
)

type Dispatcher struct {
	procs []*processor
	uid   string
}

func NewDispatcher(numOfProcs, bufferSize int) *Dispatcher {
	uuid, err := uuid.NewRandom()
	if err != nil {
		// should not panic...
		panic(err)
	}
	dispatcher := &Dispatcher{
		uid: uuid.String(),
	}

	dispatcher.procs = []*processor{}
	for i := 0; i < numOfProcs; i++ {
		dispatcher.procs = append(dispatcher.procs, NewProcessor(bufferSize))
	}

	return dispatcher
}

func (d *Dispatcher) Register(handler func(interface{})) {
	for _, proc := range d.procs {
		proc.register(handler)
	}
}

func (d *Dispatcher) Run() {
	for i := 0; i < len(d.procs); i++ {
		fmt.Printf("process %d is running\n", i)
		go d.procs[i].run()
	}
}

func (d *Dispatcher) Put(data interface{}) {
	for {
		id := rand.Intn(len(d.procs))
		fmt.Printf("process %d is executing the job\n", id)
		d.procs[id].inCh <- data
		return
	}
}

func (d *Dispatcher) Stop() {
	for i := 0; i < len(d.procs); i++ {
		d.procs[i].stop()
	}
}
