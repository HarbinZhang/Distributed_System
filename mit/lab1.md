# Lab 1: MapReduce
This lab is about building a MapReduce library with fault tolerant.

### Part I: Map/Reduce input and output
Includes: the function that divides up the output of a map task, and the function that gathers all the inputs for a reduce task. These tasks are carried out by the doMap() function in common_map.go, and the doReduce() function in common_reduce.go respectively.  

In doMap(), using hash(key) % nReduce to let messages with the same key go to the same temp file, which would be used by the same reducer.

### Part 2: Single-worker word count
Nothing special.
Just (map, reduce) function for word count.

### Part 3: Distributing MapReduce tasks
The job is to implement schedule() in mapreduce/schedule.go. The master calls schedule() twice during a MapReduce job, once for the Map phase, and once for the Reduce phase. schedule()'s job is to hand out tasks to the available workers. There will usually be more tasks than worker threads, so schedule() must give each worker a sequence of tasks, one at a time. schedule() should wait until all tasks have completed, and then return.

##### WaitGroup
- WaitGroup for waiting workers job finished.
- When to Add() or Done()? try cover as less code as possible.

##### Task transfer
Two channels:
- registerChan: input channel, for ready workers.
- readyChan: output channel, for assigning task to workers.

Master -> 
### Handling worker failures
In this part you will make the master handle failed workers. MapReduce makes this relatively easy because workers don't have persistent state. If a worker fails while handling an RPC from the master, the master's call() will eventually return false due to a timeout. In that situation, the master should re-assign the task given to the failed worker to another worker.

##### failureCounts
- sync.Mutex
- retry channel