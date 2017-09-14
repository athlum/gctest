package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"strings"
	"sync"
	"time"
)

func loopWithStats() {
	for i := 0; i < 100; i++ {
		if t := pop(); t != nil {
			t.handler()
			// t.handler = nil
		}
		var mem runtime.MemStats
		runtime.ReadMemStats(&mem)
		fmt.Printf("alloc [%v] \t heapAlloc [%v] \n", mem.Alloc, mem.HeapAlloc)

		time.Sleep(time.Second)
	}
}

type task struct {
	handler func()
	acqed   bool
}

var tasks = []*task{}

func push() {
	s := strings.Repeat("s", 1000000)
	tasks = append(tasks, &task{
		handler: func() {
			fmt.Sprintf("%vs", s)
		},
	})
}

func pop() *task {
	if len(tasks) > 0 {
		t := tasks[0]
		last := tasks[len(tasks)-1]
		tasks = tasks[1:]
		if len(tasks) > 0 {
			tasks[len(tasks)-1] = last
		}
		return t
	}
	return nil
}

func main() {
	go func() {
		http.ListenAndServe("0.0.0.0:10030", nil)
	}()
	wg := &sync.WaitGroup{}
	push()
	wg.Add(1)

	go loopWithStats()
	go func() {
		time.Sleep(time.Second * 10)
		runtime.GC()
	}()

	wg.Wait()
}
