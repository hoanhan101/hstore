package raftkv

import (
	"fmt"
)

const nservers = 5

type Command struct {
	clerk *Clerk
}

func (cmd *Command) Setup(tag string, unreliable bool, crash bool, partitions bool, maxraftstate int) *Clerk {
	cfg := makeCmdConfig(tag, nservers, unreliable, maxraftstate)
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
