package raftkv

import (
    "labrpc"
    "crypto/rand"
    "math/big"
    "sync"
)

//
//
//
type Clerk struct {
	servers []*labrpc.ClientEnd

	id        int64
	requestID int
	mu        sync.Mutex
	preLeader int
}

//
//
//
func nrand() int64 {
	max := big.NewInt(int64(1) << 62)
	bigx, _ := rand.Int(rand.Reader, max)
	x := bigx.Int64()
	return x
}

//
//
//
func MakeClerk(servers []*labrpc.ClientEnd) *Clerk {
	ck := new(Clerk)
	ck.servers = servers

	ck.id = nrand()
	ck.preLeader = 0
	ck.requestID = 0

	return ck
}

//
// fetch the current value for a key.
// returns "" if the key does not exist.
// keeps trying forever in the face of all other errors.
//
// you can send an RPC with code like this:
// ok := ck.servers[i].Call("RaftKV.Get", &args, &reply)
//
// the types of args and reply (including whether they are pointers)
// must match the declared types of the RPC handler function's
// arguments. and reply must be passed as a pointer.
//
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

//
// shared by Put and Append.
//
// you can send an RPC with code like this:
// ok := ck.servers[i].Call("RaftKV.PutAppend", &args, &reply)
//
// the types of args and reply (including whether they are pointers)
// must match the declared types of the RPC handler function's
// arguments. and reply must be passed as a pointer.
//
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

//
//
//
func (ck *Clerk) Put(key string, value string) {
	ck.PutAppend(key, value, "Put")
}

//
//
//
func (ck *Clerk) Append(key string, value string) {
	ck.PutAppend(key, value, "Append")
}
