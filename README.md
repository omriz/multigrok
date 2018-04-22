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

## Future Thoughts

+ Responses can append to the Path variable in the json. It may be expensive to Marshall and unmarahall but it will keep things consistent. 
+ Another option is to constantly wrap it in additional layers, but that can break the protocol even more. 
+ When fetching source code we can either fetch directly by decoding the path or go through the full chain. That mostly depends on the topology. In a small flat cluster the first option will work and be faster. We should consider doing both. 

### Testing

+ Create fake backend with pre-recorded responses.
+ Run the server on avaliable port. 
+ Do the cascading setup
