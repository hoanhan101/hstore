package raftkv

import (
	"fmt"
)

type Command struct {
	clerk *Clerk
}

//
// Simple setup without log compaction
//
func (cmd *Command) Setup(nservers int) *Clerk {
	cfg := makeCmdConfig(nservers, -1)
	fmt.Printf("Boot up with %v servers\n", nservers)

	ck := cfg.makeClient(cfg.All())
	cmd.clerk = ck

	return ck
}

func (cmd *Command) Put(key string, value string) string {
	cmd.clerk.Put(key, value)
	fmt.Printf("Put(%v, %v)\n", key, value)

	return key
}

func (cmd *Command) Get(key string) string {
	result := cmd.clerk.Get(key)
	fmt.Printf("Get(%v) -> %v\n", key, result)

	return result
}
