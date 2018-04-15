package mapreduce

import (
	"fmt"
	"sync"
)

//
// schedule() starts and waits for all tasks in the given phase (mapPhase
// or reducePhase). the mapFiles argument holds the names of the files that
// are the inputs to the map phase, one per map task. nReduce is the
// number of reduce tasks. the registerChan argument yields a stream
// of registered workers; each item is the worker's RPC address,
// suitable for passing to call(). registerChan will yield all
// existing registered workers (if any) and new ones as they register.
//
func schedule(
	jobName string,
	mapFiles []string,
	nReduce int,
	phase jobPhase,
	registerChan chan string,
) {
	var ntasks int
	var nOther int // number of inputs (for reduce) or outputs (for map)
	switch phase {
	case mapPhase:
		ntasks = len(mapFiles)
		nOther = nReduce
	case reducePhase:
		ntasks = nReduce
		nOther = len(mapFiles)
	}

	fmt.Printf("Schedule: %v %v tasks (%d I/Os)\n", ntasks, phase, nOther)

	// Make 10 channels of workers that receive string
	workers := make(chan string, 10)
	done := make(chan bool)

	// Define a WaitGroup
	var wg sync.WaitGroup

	// Send informations to channels concurrently
	go func() {
		for {
			select {
			case work := <-registerChan:
				workers <- work
			case <-done:
				break
			}
		}
	}()

	// Schedule tasks to free workers
	for i := 0; i < ntasks; i++ {
		select {
		case work := <-workers:
			doTaskArgs := DoTaskArgs{jobName, mapFiles[i], phase, i, nOther}
			wg.Add(1)
			var taskFunc func(string)

			// DoTask on each worker.
			// If worker fails, reassign to other free worker.
			taskFunc = func(work string) {
				if call(work, "Worker.DoTask", doTaskArgs, nil) {
					workers <- work
					wg.Done()
				} else {
					taskFunc(<-workers)
				}
			}
			go taskFunc(work)
		}
	}

	// Wait for all tasks to completed then return
	wg.Wait()
	done <- true
	fmt.Printf("Schedule: %v phase done\n", phase)
}
