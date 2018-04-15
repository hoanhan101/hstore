#!/bin/sh

# run_6824.sh - Run MIT's 6.824 tests 
# Author: Hoanh An (hoanhan@bennington.edu)
# Date: 04/15/18

export "GOPATH=$PWD"
cd "$GOPATH/src/$1"
go test
