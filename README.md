# hraft

## Usage

### CLI

`python3 mgmt.py <ARGUMENTS>` where `ARGUMENTS` are defined below:

Arguments | Description
-- | --
`-r server -c start -h <host>` | Start a seed node
`-r server -c join -h <host>` | Join a node (need at least 3 nodes for the leader election) 
`-r server -c list` | List all available nodes
`-r server -c kill -h <host>` | Kill a node 
`-r server -c stop -h <host>` | Stop a node 
`-r server -c restart -h <host>` | Restart a node 
`-r client -c read -h <host>` | Get all keys and values
`-r client -c read -h <host> -k <key>` | Read a value for a given key in a node
`-r client -c write -h <host> -k <key> -v <value>` | Write a value to a key in a node
`-r client -c update -h <host> -k <key> -v <value>` | Update a value for a key in a node
`-r client -c delete -h <host> -k <key>` | Delete a key in a node


### APIs

Endpoint | Description
-- | --
`/read` | Read all keys and values
`/read/<key>` | Read a value for a given key in a node
`/write` | Write a value to a key in a node
`/update` | Update a value for a key in a node
`/delete/<key>` | Delete a key in a node


## Testing


## Docker
