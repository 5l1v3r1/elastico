# Elastico

The **elastico** utility works on your Elasticsearch cluster(s). You can use it to gain cluster status, search through indexes and types, copy indexes from one cluster to another, create and restore backups. 

## State

Elastico is work in progress. The software is by no means ready.

## Install
Installation using Brew is most convenient. Brew will compile the source and maintain the installation.

### Install using Brew
```
$ brew tap dutchcoders/homebrew-elastico
$ brew install elastico
```

### Compiling from sourcecode

Make sure at least Go 1.6 is installed.

```
$ go get github.com/dutchcoders/elastico
```

## Contributions

Contributions are welcome.

### Creators

**Remco Verhoef**
- <https://twitter.com/remco_verhoef>
- <https://twitter.com/dutchcoders>

## Copyright and license

Code and documentation copyright 2016 Remco Verhoef.

## Usage

The `-json` parameter will give you the plain json output instead of the formatted output.

The `-host` parameter defines the url of the elasticsearch cluster to operate on.

The `-debug` parameter defines debugging mode for all requests and responses.

For convenience, the following environment variables can be set:

`ELASTICO_HOST` the server to use e.g. http://127.0.0.1:9200/

`ELASTICO_INDEX` the index to use

`ELASTICO_TYPE` the type to use

### Templates

Every output template can be customized by creating a Go template in "~/.elastico" with the name of the command. The template will get the object of the response as its parameter. For example you can create a index:stats.template to customize the output as you like.

### Cluster
#### Health
``` 
$ elastico cluster:health 
```

### State
``` 
$ elastico cluster:state
```

### Index
#### Create
Create a new index.
```
$ elastico index:create -index {name}
```

#### Delete
Delete the index.

```
$ elastico index:delete -index {name}
```

#### Copy
This command will copy an index (type) from a location to another location. This can be on the same cluster, or on a different one.  

```
$ elastico index:copy (-index {index}) (-type {type}) {dest}
```

#### Stats
This command will show you the stats of the index(es) or all indexes.

```
$ elastico index:stats (-index {index})
```

#### Get
This command will retrieve info and settings about the index(es) or all indexes.

```
$ elastico index:get (-index {index})
```

#### Recovery
```
$ elastico index:recovery -index {index} 
```

### Snapshots

#### All
Show snapshots.

```
$ elastico snapshots
```

#### Register
```
$ elastico snapshot:register -location {location} -type fs {name}
```

#### Status
```
$ elastico snapshot:status
```

### Search
Search will search through the indexes and show the relevant highlighted results.

```
$ elastico search (-disable-highlight) (-index {index}) (-type {type}) your-query
```

You can also define your own query, and pipe it to elastico.
```
$ cat query.json | elastico search (-index {index}) (-type {type})
```

### Document
#### Get
The get command allows to get a typed JSON document from the index based on its id. 

```
$ elastico document:get -index {index} -type {type} {documentid}
```

#### Put
```
$ cat doc.json | elastico document:put -index {index} -type {type} {documentid}
```

### Analyze
Performs the analysis process on a text on a specific index and return the tokens breakdown of the text.

```
$ elastico analyze -index {index} -field {field} text-to-analyze
```

### Mapping

#### Edit
Edit will open your default editor with the current mapping of the type. The modifications you'll make will be send to Elasticsearch.

```
$ elastico mapping:edit -index {index} -type {type} 
```

