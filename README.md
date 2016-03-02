# Elastico

Elastico is a console application to maintain your Elasticsearch cluster.

# State

Elastico is work in progress. The software if by no means ready.

# Build
```
$ go get github.com/constabulary/gb/...

$ gb build
```

# Usage

# Cluster
## Health
``` 
$ elastico cluster:health -host {host} 
```

## State
``` 
$ elastico cluster:health -host {host} 
```

# Index
## Create
```
$ elastico index:create {name}
```

## Delete
```
$ elastico index:delete {name}
```

## Copy
```
$ elastico index:copy -host {host} (-index {index}) (-type {type}) {dest}
```

## Stats
```
$ elastico index:stats {name}
```

## Get
```
$ elastico index:get -host {host} 
```

## Recovery
```
$ elastico index:recovery  -host {host} 
```

# Snapshots

## All
Show snapshots.

```
$ elastico snapshots
```

## Register
```
$ elastico snapshot:register -location {location} -type fs {name}
```

## Status
```
$ elastico snapshot:status
```

# Search
## Search
```
$ cat query.json | elastico --host http://127.0.0.1:9200 search -index test -type test 
```

