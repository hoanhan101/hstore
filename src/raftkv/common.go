package raftkv

// OK and ErrNoKey constants
const (
	OK       = "OK"
	ErrNoKey = "ErrNoKey"
)

// Err string type
type Err string

// PutAppendArgs structure for Put or Append Argument
type PutAppendArgs struct {
	Key       string
	Value     string
	Op        string
	ClientID  int64
	RequestID int
}

// PutAppendReply structure for Put or Append Reply
type PutAppendReply struct {
	WrongLeader bool
	Err         Err
}

// GetArgs structure for Get Argument
type GetArgs struct {
	Key       string
	ClientID  int64
	RequestID int
}

// GetReply structure for Get Reply
type GetReply struct {
	WrongLeader bool
	Err         Err
	Value       string
}
