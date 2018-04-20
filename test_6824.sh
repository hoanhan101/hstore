#!/bin/sh

#
# test_6824.sh - Run MIT's 6.824 tests for a specific package
# Author: Hoanh An (hoanhan@bennington.edu)
# Date: 04/15/18
#
# Usage:
#   ./test_6824.sh <package_name> [<test_name>]
#
# Example:
#   ./test_6824.sh kvraft: runs all tests for kvraft package
#   ./test_6824.sh raft TestInitialElection2A: runs only initial election test for raft
#

# Correct GOPATH
export "GOPATH=$PWD"
cd "$GOPATH/src/$1"

# Execute test
if [ -z  "$2" ]; then
    go test
else
    go test -run $2
fi
