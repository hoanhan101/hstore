# hstore

**hstore** is a fault-tolerant distributed key-value store. 

It is written in Go and uses the Raft consensus algorithm 
to manage a highly-available replicated log.

The project uses [MIT's 6.824: Distributed
System course](https://pdos.csail.mit.edu/6.824/) as a guideline and aims to
add more features on top of that.

Here is the project's [proposal](PROPOSAL.md). Other parts will be updated as
soon as it is live and ready.

## Testing

Assume that user set the `GOPATH` correctly, one can follow these instructions 
to run tests for different programs. If not, here is
[an example](https://github.com/hoanhan101/go-playground) on how to do it.

### MapReduce library

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

### Word Count program using MapReduce

```
$ cd hstore
$ export "GOPATH=$PWD"
$ cd "$GOPATH/src/main"
$ bash ./test-wc.sh
```

This will do the test and clean up all intermediate files afterward.
