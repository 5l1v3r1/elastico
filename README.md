# Elastico

Elastico is a console application to maintain your Elasticsearch cluster.

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
## Copy
```
$ elastico index:copy -host {host} (-index {index}) (-type {type}) {dest}
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

