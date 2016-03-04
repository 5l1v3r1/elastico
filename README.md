# Elastico

Elastico is a console application to maintain your Elasticsearch cluster. It can be used to gain information about the cluster health, but also to copy indexes from one cluster to another.

# State

Elastico is work in progress. The software if by no means ready.

# Install
```
$ brew tap dutchcoders/homebrew-elastico
$ brew install elastico
```

# Templates

The output templates can be overruled by creating a Go template file in "~/.elastico" with the name of the command. The template will get the object of the response as its parameter. For example you can create a index:stats.template to customize the output as you like.

# Build
```
$ go get github.com/constabulary/gb/...

$ gb build
```

# Usage

The `-json` parameter will give you the plain json output instead of the formatted output.

The `-host` parameter defines the host to operate on.

The `-debug` parameter defines debugging mode for all requests and responses.

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

# Get
```
$ elastico get -host {host} -index {index} -type {type} get {id}
```

# Put
```
$ cat doc.json | elastico put -host {host} -index {index} -type {type} put {id}
```

