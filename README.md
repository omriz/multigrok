# MultiGrok

## Preface
[OpenGrok](http://oracle.github.io/opengrok/) is an awesome source browser.
Using it's [docker](http://hub.docker.com/r/nagui/opengrok/) configuration it's also really easy to deploy.
However, running on a **very** large code base tends to turn problematic as it the indexer takes forever and tends to bloat.

Therefore, this project will try to implement a multiplexer over OpenGrok. It will implement the [JSON api](https://github.com/oracle/opengrok/wiki/Web-services) used by OpenGrok and will also try to forward the *xref* structure of OpenGrok to retrieve the code.

## Design Concepts

1. Stateless - the server should only hold its backends configuration and not any state of the code. Some level of caching is allowed but
you should assume it can crash at any time.
1. Concatination - you can build any topoplogy from the servers. You can build a "tree" structure of servers to do load balancing of the code searching.
1. Front-end servers will have a web interface for querying and displaying the code. They'll also implement the opengrok standard for web services *(see concatenation)*
1. *Leaf* servers can be opengrok or something else as long as they implement the same web services api. 

The server contains an [LRU Cache](https://godoc.org/github.com/hashicorp/golang-lru) that remembers succesfull direct accesses without backends.
This means that after a warmup period most of the direct links will go directly to the appropriate server (even in cascading configuration) without
overloading the network.

### Execution
For security purpuses we currently have three mode of operations:

+ **http**: Simple http server.
+ **https**: TLS encrypted communication. you need to provice cerification and key files.
+ **autoCert**: Use [Let's Encrypt](http://letsencrypt.org) to get certifications. You need to provide a host name.

### Testing

+ Create fake backend with pre-recorded responses.
+ Run the server on avaliable port. 
+ Do the cascading setup

### Service Structure
General scheme for the service:
**Frontend <-> ResponseCombiner <-> [OpenGrokBackend1, OpenGrokBackend2, SomeOtherBackend...]**
