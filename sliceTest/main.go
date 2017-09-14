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
		t := ""
		if len(data) > 0 {
			t = *data[0]
			l := data[len(data)-1]
			data[len(data)-1] = nil
			data = data[1:]
			if ll := len(data); ll > 0 {
				data[ll-1] = l
			}
		}

		var mem runtime.MemStats
		runtime.ReadMemStats(&mem)
		fmt.Printf("alloc [%v] \t heapAlloc [%v] \n", mem.Alloc, mem.HeapAlloc)

		if len(t) > 0 {
			continue
		} else {
			time.Sleep(time.Second)
		}
	}
}

var data = []*string{}

func push() {
	s := strings.Repeat("a", 1000000)
	data = append(data, &s)
}

func main() {
	push()
	go func() {
		http.ListenAndServe("0.0.0.0:10030", nil)
	}()
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go loopWithStats()
	go func() {
		time.Sleep(time.Second * 10)
		runtime.GC()
	}()

	wg.Wait()
}
