# mapreduce

## Flow

- Master starts an [RPC server](master_rpc.go) and wait for workers to 
  [register](master.go).
- As task become available, [`schedule()`](schedule.go) decides how to assign
  those tasks to workers, and how to handle worker failures.
- The master considers each input file to be one map task, and calls 
  [`doMap()`](common_map.go) at least once for each map task. It does so
  either directly (when using `Sequential()`) or by issuing the `DoTask` RPC to a
  [worker](worker.go).
  - Each call to `doMap()` reads the appropriate file, calls the map function on
    that file's contents, and writes the resulting key/value pairs to `nReduce`
    intermediate files.
  - `doMap()` hashes each key to pick the intermediate file and thus the reduce
    task that will process the key.
  - There will be `nMap` x `nReduce` files after all map tasks are done.
  - Each file name contains a prefix, the map task number, and the reduce task
    number.
  - Each worker must be able to read files written by any other worker, as well
    as the input files.
- The master next calls [`doReduce()`](common_reduce.go) at least once for each
  reduce task. As with `doMap()`, it does so either directly or through a worker.
  - The `doReduce()` for reduce task `r` collects the `r`'th intermediate file from
    each map task, and calls the reduce function for each key that appears in
    those files.
  - The reduce tasks produce `nReduce` result files.
- The master calls [`mr.merge()`](master_splitmerge.go), which merges all the
  `nReduce` files produced by the previous step into a single output.
- The master sends a Shutdown RPC to each of its workers, and then shuts down
  its own RPC server.

## Testing

```
$ cd hstore
$ export "GOPATH=$PWD" 
$ cd "$GOPATH/src/mapreduce"
$ go test -run Sequential
$ go test -run TestParallel
```

To give more verbose output, set `debugEnabled = true` in
[common.go](common.go), and add `-v` to the test command above. For example:

```
$ go test -v -run TestParallel
```
