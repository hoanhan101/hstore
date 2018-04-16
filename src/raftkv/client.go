package raftkv

import (
	"crypto/rand"
	"labrpc"
	"math/big"
	"sync"
)

// Clerk structure
type Clerk struct {
	servers []*labrpc.ClientEnd

	id        int64
	requestID int
	mu        sync.Mutex
	preLeader int
}

// nrand generates random int64 number
func nrand() int64 {
	max := big.NewInt(int64(1) << 62)
	bigx, _ := rand.Int(rand.Reader, max)
	x := bigx.Int64()
	return x
}

// MakeClerk makes a Clerk instance
func MakeClerk(servers []*labrpc.ClientEnd) *Clerk {
	ck := new(Clerk)
	ck.servers = servers

	ck.id = nrand()
	ck.preLeader = 0
	ck.requestID = 0

	return ck
}

// Get fetches the current value for a key.
// Returns "" if the key does not exist.
// Keeps trying forever in the face of all other errors.
//
// You can send an RPC with code like this:
// ok := ck.servers[i].Call("RaftKV.Get", &args, &reply)
//
// The types of args and reply (including whether they are pointers)
// Must match the declared types of the RPC handler function's
// arguments; reply must be passed as a pointer.
func (ck *Clerk) Get(key string) string {
	ck.mu.Lock()
	args := GetArgs{Key: key, ClientID: ck.id}
	args.RequestID = ck.requestID
	ck.requestID++
	ck.mu.Unlock()

	for {
		reply := GetReply{}
		ok := ck.servers[ck.preLeader].Call("RaftKV.Get", &args, &reply)
		if ok && reply.WrongLeader == false {
			return reply.Value
		}
		ck.preLeader = (ck.preLeader + 1) % len(ck.servers)
	}
}

// PutAppend is shared shared by Put and Append.
//
// You can send an RPC with code like this:
// ok := ck.servers[i].Call("RaftKV.PutAppend", &args, &reply)
//
// The types of args and reply (including whether they are pointers)
// Must match the declared types of the RPC handler function's
// arguments; reply must be passed as a pointer.
func (ck *Clerk) PutAppend(key string, value string, op string) {
	ck.mu.Lock()
	args := PutAppendArgs{Key: key, Value: value, Op: op, ClientID: ck.id}
	args.RequestID = ck.requestID
	ck.requestID++
	ck.mu.Unlock()

	for {
		reply := PutAppendReply{}
		ok := ck.servers[ck.preLeader].Call("RaftKV.PutAppend", &args, &reply)
		if ok && reply.WrongLeader == false {
			return
		}
		ck.preLeader = (ck.preLeader + 1) % len(ck.servers)
	}
}

// Put asks to put a key-value pair
func (ck *Clerk) Put(key string, value string) {
	ck.PutAppend(key, value, "Put")
}

// Append asks to append a key-value pari
func (ck *Clerk) Append(key string, value string) {
	ck.PutAppend(key, value, "Append")
}
