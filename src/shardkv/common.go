package shardkv

//
// Sharded key/value server.
// Lots of replica groups, each running op-at-a-time paxos.
// Shardmaster decides which group serves each shard.
// Shardmaster may change shard assignment from time to time.
//
// You will have to modify these definitions.
//

// Constants
const (
	OK            = "OK"
	ErrNoKey      = "ErrNoKey"
	ErrWrongGroup = "ErrWrongGroup"
)

// Err string type
type Err string

// PutAppendArgs structure
type PutAppendArgs struct {
	// You'll have to add definitions here.
	Key   string
	Value string
	Op    string // "Put" or "Append"
	// You'll have to add definitions here.
	// Field names must start with capital letters,
	// otherwise RPC will break.
}

// PutAppendReply structure
type PutAppendReply struct {
	WrongLeader bool
	Err         Err
}

// GetArgs structure
type GetArgs struct {
	Key string
	// You'll have to add definitions here.
}

// GetReply structure
type GetReply struct {
	WrongLeader bool
	Err         Err
	Value       string
}
