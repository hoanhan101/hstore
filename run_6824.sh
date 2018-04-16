#!/bin/sh

# run_6824.sh - Run MIT's 6.824 tests for a specific package
# Author: Hoanh An (hoanhan@bennington.edu)
# Date: 04/15/18
#
# Usage:
#   bash run_6824.sh <package_name>
#
# Example:
#   bash run_6824.sh kvraft will runs tests for kvraft package

export "GOPATH=$PWD"
cd "$GOPATH/src/$1"
go test
