# Note

## net/rpc or labrpc

**Problems:**

- net/prc doesn't have a Network object so cannot add/remove server
- labrpc doesn't have option for network I/O, especially Client.Dial() and Server.ServeConn()

**Solutions:**

There are 2 ways that we can go, not sure which one is more appropriate:
- Adapt laprpc to net/rpc, add more functions and rewrite the package to use websocket
- Use net/rpc and adapt labrpc library's functionalities
