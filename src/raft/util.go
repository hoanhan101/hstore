package raft

import "log"

// Debug option
const Debug = 0

// DPrintf prints debugging options
func DPrintf(format string, a ...interface{}) (n int, err error) {
	if Debug > 0 {
		log.Printf(format, a...)
	}
	return
}
