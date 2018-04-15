package raftkv

// Constants
const (
	OK       = "OK"
	ErrNoKey = "ErrNoKey"
)

// Error string
type Err string

// Put or Append structure
type PutAppendArgs struct {
	Key       string
	Value     string
	Op        string
	ClientID  int64
	RequestID int
}

// Put or Append reply structure
type PutAppendReply struct {
	WrongLeader bool
	Err         Err
}

// Get structure
type GetArgs struct {
	Key       string
	ClientID  int64
	RequestID int
}

// Get reply structure
type GetReply struct {
	WrongLeader bool
	Err         Err
	Value       string
}
