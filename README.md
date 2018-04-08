# hstore

**hstore** is a fault-tolerant distributed key-value store inspired by
[MIT's 6.824: Distributed System Spring 2017 Lab](http://nil.csail.mit.edu/6.824/2017/).
The goal of the project is to build a simple, fast and reliable database on top
of Raft, a replicated state machine protocol.

## Project Status

It is still a work-in-progress. Here is the initial project's [proposal](PROPOSAL.md).
Other parts will be updated as soon as it is live and ready.

## Table of Contents
- [Getting Started](#getting-started)
  - [Installing](#installing)
  - [Testing](#testing)
  - [Running](#running)
  - [Example](#example)
- [Reference](#reference)
- [Other](#other)
  - [mapreduce](#mapreduce)
  - [word-count](#word-count)

## Getting Started

### Installing

```
git clone https://github.com/hoanhan101/hstore.git && cd hstore
```

### Testing

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
separately, then substitute the section and run the test accordingly:
- `go test -run 2A`: leader election and heartbeats
- `go test -run 2B`: log replications
- `go test -run 2C`: persistent state

### Running

> TODO

### Example

> TODO

## Reference

- [6.824: Distributed Systems Spring 2017](http://nil.csail.mit.edu/6.824/2017/)
- [ZiyueHuang/Distributed-Systems](https://github.com/ZiyueHuang/Distributed-Systems/blob/master/src/raft/raft.go)

## Other

Other libraries/programs that are featured in MIT's lab. Here is how user can
test it.

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

### [word-count](src/main/wc.go)

```
$ cd hstore
$ export "GOPATH=$PWD"
$ cd "$GOPATH/src/main"
$ bash ./test-wc.sh
```

This will do the test and clean up all intermediate files afterward.
