# Implementing a distributed key-value store


## Introduction

By [Wikipedia](https://en.wikipedia.org/wiki/Distributed_data_store),
a distributed key-value store is a computer network where information is
stored on more than one node, often in a replicated fashion.
Examples of existing systems are 
[Apache Cassandra](http://cassandra.apache.org/), 
[Google's Bigtable](https://cloud.google.com/bigtable/), 
[Amazon's Dynamo](https://aws.amazon.com/dynamodb/),
[MongoDB](https://www.mongodb.com/), [etcd](https://coreos.com/etcd/).

It is used in production by many companies that need to solve big data
problem. It can also be found in different part of a distributed system acting
as a configuration control center. More interestingly, there are some
implementations that take the idea of a distributed key-value store and 
add more functionalities and features on top of that. One interesting
implementation is [HashiCorp's Consul](https://www.consul.io/) where 
they build a distributed, highly available, and data center aware solution 
to connect and configure applications across dynamic, distributed infrastructure. 
There is also [Redis](https://redis.io/) that is an in-memory data structure store, 
used as a database, cache and message broker.

My goal for this project is to be able to implement a distributed key-value store
from scratch as well as to gain a better understanding of some specific topics in
distributed systems, such as: consensus algorithm, distributed hash table,
RPC, CAP theorem. Another goal is to be familiar building distributed systems 
using [Go](https://golang.org/) language.

In many sections below, I will talk about the design/architecture of the
project, final product, testing. I will also provide a timeline 
with tasks, approaches and deliverables in a bi-weekly basis.


## Design

There are 3 big components of the system: a key-value storage, consensus
algorithm and routing/service discovery.

### A key-value storage

The most straightforward way is to use a hash table to store key-value pairs.
It allows user to read and write in constant time. It's also very easy to use.
In modern programming languages, hash table data structure is normally built-in
as a form of map (or dictionary) with supported CRUD (Create, Read, Update,
Delete) operations.

However, using a hash table means that I need to store everything in memory, which
is not great when the data get big. One way to solve it is to store them in
disk and use a cache system (something like LRU). Frequently visited data is kept in
memory and the rest is on disk.

> There might be other caching systems to learn from such as Redis and Memcached.

### Consensus algorithm

A consensus algorithm is critical in a distributed system because it allows a
collection of machines to work as a coherent group that can survive failures of
some its members. Paxos is undeniably the most popular algorithm. However, it is also
known for its complexity. Understanding Paxos is hard. Last term I attempted to
implement Paxos but it did not turn out very well. There were still a lot of 
aspects that I was uncertain about. Therefore, I want to try something new this
term.

[Raft](https://raft.github.io/) seems like a good fit. 
It is made to solve Paxos's understandability problem. It has been used by a 
lot of organizations, such as etcd, HashiCorp's Consul, Docker Swarm,... 
and continued to gain its popularity.

> HashiCorp even provides a nice implementation of the Raft and it is
> imported by [many systems](https://godoc.org/github.com/hashicorp/raft?importers).

### Routing/Service Discovery

The last piece of the system is routing/service discovery. At the moment,
I am not sure how to do this yet. I know that HashiCorp's Consul achieve this 
by using a DNS routing mechanism but I am not familiar with its implementation.
In the [final product section](#final-product), I will give an example of
a service discovery's behavior that I want.

> **TODO:** Spend more time looking at HashiCorp's Consul documentation.


## Flow

> There should be some images in this section for the reader to visualize the
> system easily. It also helps forming the flow of the system. I don't have one
> yet and I can't think of any at the moment. I will update this as I take a 
> closer look at Raft as well as other documents.

After looking at some similar systems, I realize that getting the consensus 
algorithm right in the first place is the most important job.
As long as I have all the nodes perform resiliently using the protocol, 
building a key-value store on top seems much more natural.
In other word, Raft does most of the heavy lifting for the system.

### Step by step

- I first start with a seed node as a server.
- I use other node to join the seed node or I can choose to join other node
  that are available in the system if I know its host. However, if I am able to
  implement the service discovery feature, I just need to tell it to `join` 
  without explicitly specifying the host. It automatically know how to route 
  to the right cluster. 

  > I think this is how a service discovery should work. However I am not sure.

- When there are 3 nodes, leader election will occur. One is the leader, the
  rest are followers.
- Using the client machine, I can read, write, update or delete key-values in any
  of these nodes.
- The key-values should be replicated among themselves, no matter where I put
  them. If I put a key-value pair in the master node then it will eventually
  propagated to other nodes. However, if I put a key-value pair in the
  non-master node, it should redirect the request to the master and do the
  consensus checking there. 

  > In this situation, is it better to put a load balancer in front of
  > these non-master nodes?

- If I choose to stop a node or multiple nodes, the system must still work.
- If I stop the master, it will start the leader election again if and only if
  the quorum size is big enough. Everything should behave the same way.


## Timeline

I am using [MIT's 6.824 Distributed System
Spring 2017](http://nil.csail.mit.edu/6.824/2017/) 
as a guideline for my implementation. Their labs provide pointers and 
resources on how to build a fault-tolerant key-value storage, which is exactly
what I want to accomplish for this project. They start with Raft implementation, 
add a distributed key-value service on top of it and eventually exploring 
the idea of sharding.

My timeline is also largely dependent on their course's schedule.

### Week 1-2
- Task:
  - Finish first draft of the proposal.
- Approach:
  - Read about similar systems and learn how do they implement it.
  - Come up with a solution myself that fits the scope of the project.
- Deliverables:
  - A reasonable well-written first draft.

### Week 3-4
- Task:
  - Get familiar with Go by going through MapReduce's implementation.
  - Start thinking about Raft implementation and update the proposal along the way.
- Approach:
  - [Lab 1: MapReduce](http://nil.csail.mit.edu/6.824/2017/labs/lab-1.html)
- Deliverables:
  - A working MapReduce library.

### Week 5-6
- Task:
  - Implement a minimum version of Raft.
- Approach:
  - [Lab 2: Raft](http://nil.csail.mit.edu/6.824/2017/labs/lab-raft.html)
- Deliverables:
  - A working version of Raft.

### Week 7-8
- Task:
  - Build a key-value store using Raft library.
- Approach:
  - [Lab 3: Fault-tolerant Key/Value
    Service](http://nil.csail.mit.edu/6.824/2017/labs/lab-kvraft.html)
  - [Lab 4: Sharded Key/Value
    Service](http://nil.csail.mit.edu/6.824/2017/labs/lab-shard.html)
- Deliverables:
  - A robust key-value store that "shards

> For now, I am leaving the rest of the weeks' tasks as *TODO*. Will update
> when I have a better understanding of the system.

### Week 9-10
- Task:
  - Implement a service discovery feature
- Approach:
  - TODO
- Deliverables:
  - TODO

### Week 11-12
- Task:
- Approach:
- Deliverables:

### Week 13-14
- Task:
- Approach:
- Deliverables:


## Final Product

This is how I see it working as the final product.

### CLI

```
NAME:
   hstore - hstore shell

USAGE:
   hstore [global options] role [role options] command [command options] [-h <host>] [-k <key>] [-v <value>]
VERSION:
   0.1.0

AUTHORS:
   Hoanh An <hoanhan@bennington.edu>

COMMANDS:
     start          Start a node
     join           Join a node to another node
     list           List all avaiable nodes
     kill           Kill a node
     stop           Stop a node
     restart        Restart a node 
     read           Read a value for a key 
     write          Write a value to a key
     update         Update a value for a key 
     delete         Delete a key-value pair 
     help, h        Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h                                  show help
   --version, -v                               print the version
```

Here are the lists of example commands: 

Commands | Description
-- | --
`hstore server start -h <host>` | Start a seed node with a given host and prompt user into the shell.
`hstore server join [-h <host>]` | Join a node to the cluster and prompt user into the shell.
`hstore server list` | List all available nodes showing their name, address, health status and type.
`hstore server kill -h <host>` | Kill a node with a given host.
`hstore server stop -h <host>` | Stop a node with a given host.
`hstore server restart -h <host>` | Restart a node with a given host.
`hstore client read -h <host>` | Get all the keys and values for a given host.
`hstore client read -h <host> -k <key>` | Read a value for a given key, for a given host.
`hstore client write -h <host> -k <key> -v <value>` | Write a value to a key for a given host.
`hstore client update -h <host> -k <key> -v <value>` | Update a value for a key for a given host.
`hstore client delete -h <host> -k <key>` | Delete a key for a given host.

### APIs

List of exposed APIs for each node.

Method | Endpoint | Description
-- | -- | --
`GET` | `/read` | Read all keys and values.
`GET` | `/read/<key>` | Read a value for a given key in a node.
`POST` | `/write` | Write a value to a key in a node.
`POST` | `/update` | Update a value for a key in a node.
`GET` | `/delete/<key>` | Delete a key in a node.


## Testing

> **Idea:** Unit test, integration test, end-to-end test 

> Other than unit test, integration test (if needed), end to end test (if
> needed), how to introduce failure injections/exercises for the system,
> exploring its behavior in the face of crashes and network partitioning?

> If the system fails in such a way that it can not function properly anymore,
  how would I recover/bring everything back up gracefully?


## Report

> **Idea:** A dashboard and a write-up paper with discussion.

> Gather data and do different types of analysis for the system here


## UI

> **Idea:** A interactive webpage.

> Something like [etcd's playground](http://play.etcd.io/play) is nice to have
> for better visualization.

## Future features

- Full text search like Elasticsearch
