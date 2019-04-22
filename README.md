# Redis Dictator

[![Go Report Card](https://goreportcard.com/badge/github.com/Junonogis/redis-dictator)](https://goreportcard.com/report/github.com/Junonogis/redis-dictator)
[![Build Status](https://travis-ci.org/Junonogis/redis-dictator.svg?branch=master)](https://travis-ci.org/Junonogis/masterservice-operator)

[Docker Repository](https://hub.docker.com/r/junonogis/redis-dictator)

Redis Dictator is a simple program running next to Redis Server in a Kubernetes deployment. His role is to provide replication and promote API over a Redis node.

Before kubernetes, dictator directly handled the master election process in Zookeeper. Now, this process is done by a kubernetes operator.
    
## Run

### Command line
```
dictator --config /etc/dictator/dictator.conf
```

### Docker

You can find a helm Chart in repository using [MasterService Operator](https://github.com/Junonogis/masterservice-operator) as a controller over Redis cluster.

Some E2E testing is done in the operator pipeline.

## Configuration
    
The Dictator configuration is a JSON file located at `/etc/dictator/dictator.json`

    $ cat /etc/dictator/dictator.conf
    {
        "svc_name" : "default",
        "log_level" : "INFO",
        "http_port": 8000,
        "zk_hosts": ["127.0.0.1"],
        "node" : {
            "name" : "local",
            "host" : "127.0.0.1",
            "port" : 6379,
            "loading_timeout" : 30
        } 
    }

The main section is composed by:
 
- `svc_name`: The Service/Cluster Name (default is `local`)
- `log_level`: The log level `DEBUG`, `INFO`, `WARN`, `FATAL` (default is `INFO`)
- `http_port`: The port of the HTTP listener, used for interact with Dictator (Disable/Enable)
- `node`: The Redis node info (detailed bellow)
- `master_service`: Name of the master kubernetes service

The node section is composed by:

- `name`: The server name, FQDN or "display" name (default is `local`)
- `host`: The Address of the Redis server (default is `localhost`)
- `port`: The Port of the Redis server (default is `6379`)
- `loading_timeout`: The time in second that Dictator accepts to wait during Redis loads its dataset to memory

## Redis at BlaBlaCar

### Usages

We are using Redis at BlaBlaCar since some couples of years, for a large spectrum of usecases:
- to cache data (obvious no?)
- maintain counters (ex: "quota" usages)
- create functional locks
- stored some configuration/feature swicthes
- ...
 
Basically, our use cases involved that we should provide a "quite" high available solution, the Redis replication and persistence solutions allow us to propose a satisfying HA. Provided that a master/slaves topology should be quickly/well reconfigured in case of master failure...

### Motivation
We spent lot of time/energy to test some HA & Cluster solutions around Redis. We put aside the idea of clusterize Redis (in term of turnkey sharding solution), it complexify your topologies and reduce drastically your consistency... We choosed to shard our dataset functionally by creating several master/slaves clusters instead of one magical auto-sharding, auto-scaling, auto-[...]ing black box.

We are not the first to develop tooling to manipulate master/slaves toplologies. The most known is surely [Redis Sentinel](http://redis.io/topics/sentinel) but the configuration file rewriting bother us a little (note that we are in a full containers context at BlaBlaCar).  By the way, we should admit that our main motivation is certainly because developping your own tool is fun and offers a lot of advantages, the solution fits perfectly to your needs, you can chose the language, merge PRs quickly...