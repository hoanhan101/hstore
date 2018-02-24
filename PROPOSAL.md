# Implementing a distributed key/value store

In this paper, I will describe my plan to build a distributed key/value store.
The goal of this project to gain a better understanding of some specific topics
in distributed systems, such as: consensus protocol, RPC, CAP theorem.

There are 3 big components that I think it will have. 

The first one is obviously the key/value store. I am thinking about using
memory to store all the keys and values. Another option can be building a
Distributed Hash Table.

The second one is consensus protocol. I am thinking about using Raft, instead
of Paxos. The main reason is because of understandability that Raft has. Last
term I attemped to implement Paxos but it did not turn out very well. There
were still a lot of aspects that I was uncertain. Hopefully Raft is easier.

The third one is routing/service discovery. At the momemnt, I am not sure how
to do this yet. I know that Consul achieve this by DNS routing but that's all I
know.

In the sections below, I will provide more details about the
design/architecture and outline of the project.

# Design

TODO

# Timeline

TODO

Here is the format that every week should follow:
- Task
- Approach
- Deliverables

# Final Product

Here is how I see it working as the final product.

CLI looks like this:
```
Usage: python3 mgmt.py -r <role> -c <command> [-h <host>] [-k <key>] [-v
<value>]
Help:  Run the management system to control and execute tasks.
Options:
      -r                Role to be defined.
      -c                Command to be executed.
      -h                Host name, including address and port.
      -k                Key to be stored.
      -v                Value corresponding to the key.
```

Here are the lists of arguments with descriptions:

Arguments | Description
-- | --
`-r server -c start -h <host>` | Start a seed node
`-r server -c join -h <host>` | Join a node
`-r server -c list` | List all available nodes
`-r server -c kill -h <host>` | Kill a node
`-r server -c stop -h <host>` | Stop a node
`-r server -c restart -h <host>` | Restart a node
`-r client -c read -h <host>` | Get all the keys and values
`-r client -c read -h <host> -k <key>` | Read a value for a given key
`-r client -c write -h <host> -k <key> -v <value>` | Write a value to a key
`-r client -c update -h <host> -k <key> -v <value>` | Update a value for a key
`-r client -c delete -h <host> -k <key>` | Delete a key

You first start with a seed node.
You join the seed node with as many as you want.
When there are 3 nodes, leader election will occur.
One is the leader, the rest are followers.
As a client, you can read, write, update, delete a key/value in any of these
nodes.
The key/value should be replicated among themselves, no matter where you put it
(it doesn't have to be the master)
If you stop a node or mutiple nodes, the system should still work.
If you stop the master, it will start the leader election again and everything
should still work.

# Testing

There should be automated tests for the system.

# Docker

Everything should be dockerized.
