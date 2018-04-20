# Note

## net/rpc or labrpc

**Problems:**

- net/prc doesn't have a Network object so cannot add/remove server
- labrpc doesn't have option for network I/O (e.g.: Client.Dial and Server.ServeConn)

**Solutions:**

There are 2 ways that we can go, not sure which one is more appropriate:
- Adapt laprpc to net/rpc, add more functions and rewrite the package to use websocket
- Use net/rpc and adapt labrpc library's functionalities

Other thoughts:
- How to plug in the networking part without rewriting everything there?
- How to make each server in the quorum expose a port and still keep the logic of RaftKV behind
  the scene? Because RPC call can be achieve through Go routines, all we need is a network
  interface. 
- In labrpc, why does a client need to connect to every server?
