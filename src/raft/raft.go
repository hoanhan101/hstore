package raft

//
// Reference: https://github.com/xingdl2007/6.824-2017/tree/master/src/raft
//
// This is an outline of the API that raft must expose to
// the service (or tester). see comments below for
// each of these functions for more details.
//
// rf = Make(...)
//   Create a new Raft server.
// rf.Start(command interface{}) (index, term, isleader)
//   Start agreement on a new log entry
// rf.GetState() (term, isLeader)
//   Ask a Raft for its current term, and whether it thinks it is leader
// ApplyMsg
//   Each time a new entry is committed to the log, each Raft peer
//   should send an ApplyMsg to the service (or tester)
//   in the same server.
//

import (
	"labrpc"
	"math/rand"
	"sync"
	"time"
)

// import "bytes"
// import "labgob"

//
// As each Raft peer becomes aware that successive log entries are
// committed, the peer should send an ApplyMsg to the service (or
// tester) on the same server, via the applyCh passed to Make(). set
// CommandValid to true to indicate that the ApplyMsg contains a newly
// committed log entry.
//
// In Lab 3 you'll want to send other kinds of messages (e.g.,
// snapshots) on the applyCh; at that point you can add fields to
// ApplyMsg, but set CommandValid to false for these other uses.
//
type ApplyMsg struct {
	CommandValid bool
	Command      interface{}
	CommandIndex int
}

//
// Log Entry
//
type LogEntry struct {
	Term    int
	Command interface{}
}

//
// A Go object implementing a single Raft peer.
//
type Raft struct {
	mu        sync.Mutex          // Lock to protect shared access to this peer's state
	peers     []*labrpc.ClientEnd // RPC end points of all peers
	persister *Persister          // Object to hold this peer's persisted state
	me        int                 // this peer's index into peers[]

	// Election Term
	isLeader bool

	resetTimer      chan struct{}
	electionTimer   *time.Timer
	electionTimeout time.Duration // 500 ~ 800 ms

	// Heartbeat
	heartbeatInterval time.Duration // 150 ms

	// Persistent state on all servers
	// Update on stable storage before responding to RPCs
	currentTerm int
	votedFor    int
	logs        []LogEntry

	// Volatile state on all servers
	commitIndex int
	lastApplied int

	// Volatile state on leaders
	// Reinitialized after election
	nextIndex  []int
	matchIndex []int
}

//
// Return currentTerm and whether this server
// believes it is the leader.
//
func (rf *Raft) GetState() (int, bool) {
	var term int
	var isleader bool

	rf.mu.Lock()
	defer rf.mu.Unlock()

	term = rf.currentTerm
	isleader = rf.isLeader

	return term, isleader
}

//
// Save Raft's persistent state to stable storage,
// where it can later be retrieved after a crash and restart.
// See paper's Figure 2 for a description of what should be persistent.
//
func (rf *Raft) persist() {
	// Your code here (2C).
	// Example:
	// w := new(bytes.Buffer)
	// e := labgob.NewEncoder(w)
	// e.Encode(rf.xxx)
	// e.Encode(rf.yyy)
	// data := w.Bytes()
	// rf.persister.SaveRaftState(data)
}

//
// Restore previously persisted state.
//
func (rf *Raft) readPersist(data []byte) {
	if data == nil || len(data) < 1 { // bootstrap without any state?
		return
	}
	// Your code here (2C).
	// Example:
	// r := bytes.NewBuffer(data)
	// d := labgob.NewDecoder(r)
	// var xxx
	// var yyy
	// if d.Decode(&xxx) != nil ||
	//    d.Decode(&yyy) != nil {
	//   error...
	// } else {
	//   rf.xxx = xxx
	//   rf.yyy = yyy
	// }
}

//
// RequestVote RPC arguments structure.
// Field names must start with capital letters!
//
type RequestVoteArgs struct {
	Term         int // candidate's term
	CandidateID  int // candidate's requesting vote
	LastLogIndex int // index of candidate's last log entry
	LastLogTerm  int // term of candidate's last log entry
}

//
// Initialize RequestVote RPC arguments
//
func (rf *Raft) initRequestVoteArgs(args *RequestVoteArgs) {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	// Turn to candidate and vote to itself
	rf.votedFor = rf.me
	rf.currentTerm += 1

	args.Term = rf.currentTerm
	args.CandidateID = rf.me
	args.LastLogIndex = len(rf.logs) - 1
	args.LastLogTerm = rf.logs[args.LastLogIndex].Term
}

//
// RequestVote RPC reply structure.
// Field names must start with capital letters!
//
type RequestVoteReply struct {
	Term        int  // currentTerm, for candidate to update itself
	VoteGranted bool // true means candidate received vote
}

//
// RequestVote RPC handler.
//
func (rf *Raft) RequestVote(args *RequestVoteArgs, reply *RequestVoteReply) {
	// Your code here (2A, 2B).
	rf.mu.Lock()
	defer rf.mu.Unlock()

	// If one's current term is smaller than other's,
	// then update its current term to the larget value.
	// Otherwise, its term is out of date. Revert to follower state.
	if args.Term < rf.currentTerm {
		reply.Term = rf.currentTerm
		reply.VoteGranted = false
	} else {
		// Convert to follower
		if args.Term > rf.currentTerm {
			rf.currentTerm = args.Term
			rf.votedFor = -1
		}

		// If voteFor is null or candidateId (candidate itself),
		// and candidate's log is at least up-to-date as receiver's log
		// then grant vote
		if rf.votedFor == -1 || rf.votedFor == args.CandidateID {
			lastLogIndex := len(rf.logs) - 1
			lastLogTerm := rf.logs[lastLogIndex].Term

			if (args.LastLogIndex == lastLogTerm && args.LastLogIndex >= lastLogIndex) ||
				args.LastLogTerm > lastLogTerm {
				rf.isLeader = false
				rf.votedFor = args.CandidateID
				reply.VoteGranted = true

				// Reset timer after granting a vote to another peer
				rf.resetTimer <- struct{}{}
			}
		}
	}
}

//
// Example code to send a RequestVote RPC to a server.
// Server is the index of the target server in rf.peers[].
// Expects RPC arguments in args.
// Fills in *reply with RPC reply, so caller should
// pass &reply.
// The types of the args and reply passed to Call() must be
// the same as the types of the arguments declared in the
// handler function (including whether they are pointers).
//
// The labrpc package simulates a lossy network, in which servers
// may be unreachable, and in which requests and replies may be lost.
// Call() sends a request and waits for a reply. If a reply arrives
// within a timeout interval, Call() returns true; otherwise
// Call() returns false. Thus Call() may not return for a while.
// A false return can be caused by a dead server, a live server that
// can't be reached, a lost request, or a lost reply.
//
// Call() is guaranteed to return (perhaps after a delay) *except* if the
// handler function on the server side does not return.  Thus there
// is no need to implement your own timeouts around Call().
//
// Look at the comments in ../labrpc/labrpc.go for more details.
//
// If you're having trouble getting RPC to work, check that you've
// capitalized all field names in structs passed over RPC, and
// that the caller passes the address of the reply struct with &, not
// the struct itself.
//
func (rf *Raft) sendRequestVote(server int, args *RequestVoteArgs, reply *RequestVoteReply) bool {
	ok := rf.peers[server].Call("Raft.RequestVote", args, reply)
	return ok
}

//
// AppendEntries RPC arguments structure
// Invoked by leader to replicate log entries (section 5.3);
// also used as heartbeat (section 5.2)
//
type AppendEntriesArgs struct {
	Term         int        // leader's term
	LeaderID     int        // so follower can redirect clients
	PrevLogIndex int        // index of log entry immediately preceding new ones
	PrevLogTerm  int        // term of prevLogIndex entry
	Entries      []LogEntry // log entries to store (empty for heartbeat; may send more than one for eddiciency)
	LeaderCommit int        // leader's commitIndex
}

//
// Initialize AppendEntries RPC arguments
//
func (rf *Raft) initAppendEntriesArgs(args *AppendEntriesArgs, heartbeat bool) {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	args.Term = rf.currentTerm
	args.LeaderID = rf.me
	args.PrevLogIndex = len(rf.logs) - 1
	args.PrevLogTerm = rf.logs[args.PrevLogIndex].Term
	args.LeaderCommit = rf.commitIndex

	if heartbeat {
		args.Entries = nil
	}
}

//
// AppendEntries RPC reply structure
//
type AppendEntriesReply struct {
	Term    int  // currentTerm, for leader to update itself
	Success bool // true if follower contained entry matching prevLogIndex and prevLogTerm
}

//
// AppendEntries RPC Handler
//
func (rf *Raft) AppendEntries(args *AppendEntriesArgs, reply *AppendEntriesReply) {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	// If heartbeat message, suppress leader election
	if len(args.Entries) == 0 {
		if args.Term < rf.currentTerm {
			reply.Term = rf.currentTerm
			reply.Success = false
		} else {
			reply.Success = true

			// If encounter a leader, change to follower
			if rf.isLeader {
				rf.isLeader = false
			}

			// Reset election timer
			rf.resetTimer <- struct{}{}
		}
	}
}

//
// Send a AppendEntries RPC to a server
// Model after sendRequestVote
//
func (rf *Raft) sendAppendEntries(server int, args *AppendEntriesArgs, reply *AppendEntriesReply) bool {
	ok := rf.peers[server].Call("Raft.AppendEntries", args, reply)
	return ok
}

//
// The service using Raft (e.g. a k/v server) wants to start
// agreement on the next command to be appended to Raft's log. If this
// server isn't the leader, returns false. Otherwise start the
// agreement and return immediately. There is no guarantee that this
// command will ever be committed to the Raft log, since the leader
// may fail or lose an election.
//
// The first return value is the index that the command will appear at
// if it's ever committed. The second return value is the current
// term. The third return value is true if this server believes it is
// the leader.
//
func (rf *Raft) Start(command interface{}) (int, int, bool) {
	index := -1
	term := -1
	isLeader := true

	// Your code here (2B).

	return index, term, isLeader
}

//
// The tester calls Kill() when a Raft instance won't
// be needed again. You are not required to do anything
// in Kill(), but it might be convenient to (for example)
// turn off debug output from this instance.
//
func (rf *Raft) Kill() {
	// Your code here, if desired.
}

//
// The service or tester wants to create a Raft server. the ports
// of all the Raft servers (including this one) are in peers[]. this
// server's port is peers[me]. All the servers' peers[] arrays
// have the same order. Persister is a place for this server to
// save its persistent state, and also initially holds the most
// recent saved state, if any. applyCh is a channel on which the
// tester or service expects Raft to send ApplyMsg messages.
// Make() must return quickly, so it should start goroutines
// for any long-running work.
//
func Make(peers []*labrpc.ClientEnd, me int,
	persister *Persister, applyCh chan ApplyMsg) *Raft {
	rf := &Raft{}
	rf.peers = peers
	rf.persister = persister
	rf.me = me

	rf.isLeader = false
	rf.votedFor = -1

	// Reinitialzie after election
	rf.nextIndex = make([]int, len(peers))
	rf.matchIndex = make([]int, len(peers))

	// Start with logs of length 1
	rf.logs = make([]LogEntry, 1)

	// matchIndex is initialized to leader's lastLogIndex + 1
	for i := 0; i < len(peers); i++ {
		rf.matchIndex[i] = len(rf.logs)
	}

	// Set election timeout to 400-800 ms
	rf.electionTimeout = time.Millisecond * (400 + time.Duration(rand.Int63()%400))
	rf.electionTimer = time.NewTimer(rf.electionTimeout)
	rf.resetTimer = make(chan struct{})

	// Set heartbeat interval to 100 ms
	rf.heartbeatInterval = time.Millisecond * 100

	// Debugging message
	DPrintf("Peer %d : Election(%s) Heartbeat(%s)\n", rf.me, rf.electionTimeout, rf.heartbeatInterval)

	// Initialize from state persisted before a crash
	rf.readPersist(persister.ReadRaftState())

	// Kick off election
	go rf.electionDaemon()

	return rf
}

//
// Election Daemon never exit, but can be reset.
// It listens in the background and reset timer based on different cases.
//
func (rf *Raft) electionDaemon() {
	for {
		select {
		case <-rf.resetTimer:
			rf.electionTimer.Reset(rf.electionTimeout)
			// When the Timer expires, the current time will be sent on C
		case <-rf.electionTimer.C:
			rf.electionTimer.Reset(rf.electionTimeout)
			go rf.canvassVotes()
		}
	}
}

//
// CanvassVotes issues RequestVote RPC - Election logic
//
func (rf *Raft) canvassVotes() {
	var voteArgs RequestVoteArgs
	rf.initRequestVoteArgs(&voteArgs)

	// Use buffered channel to avoid goroutine leak
	peers := len(rf.peers)

	// Make channel of RequestVoteReply?
	replies := make(chan RequestVoteReply, peers)

	// Send RequestVoteRPC to all other servers
	var wg sync.WaitGroup
	for i := 0; i < peers; i++ {
		if i == rf.me {
			// Reset itself electionTimer
			rf.resetTimer <- struct{}{}
		} else {
			wg.Add(1)
			go func(n int) {
				defer wg.Done()

				var reply RequestVoteReply
				if rf.sendRequestVote(n, &voteArgs, &reply) {
					replies <- reply
				}
			}(i)
		}
	}

	// Remember to close the channel
	go func() {
		wg.Wait()
		close(replies)
	}()

	var votes = 1
	for reply := range replies {
		if reply.VoteGranted == true {
			// If receives votes from majority of servers, becomes a leader
			if votes++; votes > peers/2 {
				rf.mu.Lock()
				rf.isLeader = true
				rf.mu.Unlock()

				// Send AppendEntries heartbeats to all servers
				go rf.heartbeatDaemon()
				return
			}
			// If we found a new leader, step down to be a follower
		} else if reply.Term > voteArgs.Term {
			rf.mu.Lock()
			rf.isLeader = false
			rf.votedFor = -1
			rf.mu.Unlock()
			return
		}
	}
}

//
// Heartbeat Daemon will exist when peer is not a leader anymore
//
func (rf *Raft) heartbeatDaemon() {
	for {
		if _, isLeader := rf.GetState(); isLeader {
			for i := 0; i < len(rf.peers); i++ {
				if i == rf.me {
					// Reset itself electionTimer
					rf.resetTimer <- struct{}{}
				} else {
					go rf.heartbeat(i)
				}
			}
			time.Sleep(rf.heartbeatInterval)
		} else {
			break
		}
	}
}

//
// Process replay of AppendEntries including heartbeat
//
func (rf *Raft) heartbeat(n int) {
	var args AppendEntriesArgs
	rf.initAppendEntriesArgs(&args, true)

	var reply AppendEntriesReply
	if rf.sendAppendEntries(n, &args, &reply) {
		// Retry until it succeed?
		if !reply.Success {
			rf.mu.Lock()
			defer rf.mu.Unlock()

			// If we found a new leader, step down to be a follower
			if reply.Term > rf.currentTerm {
				rf.currentTerm = reply.Term
				rf.isLeader = false
				rf.votedFor = -1
			} else {
				// TODO: Send log entry from leader
			}
		}
	}
}
