# raft

## Outline API

- A service calls `Make(peers,me,…)` to create a Raft peer.
  - The peers argument is an array of network identifiers of 
  the Raft peers (including this one), for use with labrpc RPC.
  - The `me` argument is the index of this peer in the peers array.
- `Start(command)` asks Raft to start the processing to append the command to
the replicated log.
  - `Start()` should return immediately, without waiting for the log appends to
    complete.
- The service expects your implementation to send an `ApplyMsg` for each newly
  committed log entry to the `applyCh` argument to `Make()`
