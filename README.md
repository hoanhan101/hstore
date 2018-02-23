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
`-r client -c read -h <host>` | Get all keys/values
`-r client -c read -h <host> -k <key> -v <value>` | Read a key/value for a given host
`-r client -c write -h <host> -k <key> -v <value>` | Write a key/value for a given host
`-r client -c update -h <host> -k <key> -v <value>` | Update a key/value for a given host
`-r client -c delete -h <host> -k <key> -v <value>` | Delete a key/value for a given host


### APIs

Endpoint | Description
-- | --
`/read` | Read all keys/values
`/read/<key>` | Read a key value for a given host
`/write` | Write a key value for a given host
`/update` | Update a key value for a given host
`/delete/<key>` | Delete a key value for a given host


## Testing


## Docker
