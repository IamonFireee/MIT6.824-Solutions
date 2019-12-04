package mapreduce

import (
	"fmt"
)

//
// schedule() starts and waits for all tasks in the given phase (Map
// or Reduce). the mapFiles argument holds the names of the files that
// are the inputs to the map phase, one per map task. nReduce is the
// number of reduce tasks. the registerChan argument yields a stream
// of registered workers; each item is the worker's RPC address,
// suitable for passing to call(). registerChan will yield all
// existing registered workers (if any) and new ones as they register.
//
func schedule(jobName string, mapFiles []string, nReduce int, phase jobPhase, registerChan chan string) {
	var ntasks int
	var n_other int // number of inputs (for reduce) or outputs (for map)
	switch phase {
	case mapPhase:
		ntasks = len(mapFiles)
		n_other = nReduce
	case reducePhase:
		ntasks = nReduce
		n_other = len(mapFiles)
	}

	fmt.Printf("Schedule: %v %v tasks (%d I/Os)\n", ntasks, phase, n_other)

	// construct the DoTaskArgs for call()
	arg:=func(task int) DoTaskArgs{
		var arg DoTaskArgs
		arg.Phase=phase
		arg.TaskNumber=task
		arg.JobName=jobName
		arg.NumOtherPhase=n_other
		if phase==mapPhase{
			arg.File=mapFiles[task]
		}
		return arg
	}

	// create a uncompleted tasks queue
	tasks:=make(chan int)
	go func(){
		for i:=0;i<ntasks;i++{
			tasks<-i // put all tasks into the queue
		}
	}()
	successTasks:=0
	successChan:=make(chan int) // for threads tell us if their task finished
	loop:
	for{
		select {
		case task:=<-tasks:
			go func() {
				worker := <-registerChan
				status := call(worker, "Worker.DoTask", arg(task), nil)
				if status {
					successChan <- task
					registerChan<-worker
				} else {
					tasks <- task
				}
			}()
		case <-successChan:
			successTasks+=1 // another task finished
		default:
			if successTasks==ntasks{ // all task finished
				break loop
			}
		}
	}

	fmt.Printf("Schedule: %v phase done\n", phase)
}

