# hstore

**hstore** is a fault-tolerant distributed key-value store. 

It is written in Go and uses the Raft consensus algorithm 
to manage a highly-available replicated log.

The project uses [MIT's 6.824: Distributed
System Spring 2017](http://nil.csail.mit.edu/6.824/2017/) as a guideline
and foundation to add more interesting features.

Here is the project's [proposal](PROPOSAL.md). Other parts will be updated as
soon as it is live and ready.

## Testing

Assume that user set the `GOPATH` correctly, one can follow these instructions 
to run tests for different programs. If not, here is
[an example](https://github.com/hoanhan101/go-playground) on how to do it.

### [raft](src/raft)

```
$ cd hstore
$ export "GOPATH=$PWD" 
$ cd "$GOPATH/src/raft"
$ go test
```

This will run all the test for Raft. If one want to test the program
separately, then substitute the section and run the test accordingly.
For example, in order to test leader election and heartbeats, one can run `go
test -run 2A` since

> The goal for Part 2A is for a single leader to be elected, for the leader to
> remain the leader if there are no failures, and for a new leader to take over
> if the old leader fails or if packets to/from the old leader are lost.

Part 2B and 2C are work-in-progress.

### [mapreduce](src/mapreduce)

```
$ cd hstore
$ export "GOPATH=$PWD" 
$ cd "$GOPATH/src/mapreduce"
$ go test -run Sequential
$ go test -run TestBasic
```

To give more verbose output, set `debugEnabled = true` in
[common.go](src/mapreduce/common.go), and add `-v` to the test command above. 
For example:

```
$ go test -v -run TestBasic
```

### [word-count with mapreduce](src/main/wc.go)

```
$ cd hstore
$ export "GOPATH=$PWD"
$ cd "$GOPATH/src/main"
$ bash ./test-wc.sh
```

This will do the test and clean up all intermediate files afterward.
