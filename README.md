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

