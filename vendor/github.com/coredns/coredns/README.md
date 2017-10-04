[![CoreDNS](https://coredns.io/images/CoreDNS_Colour_Horizontal.png)](https://coredns.io)

[![Documentation](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/coredns/coredns)
[![Build Status](https://img.shields.io/travis/coredns/coredns/master.svg?label=build)](https://travis-ci.org/coredns/coredns)
[![Code Coverage](https://img.shields.io/codecov/c/github/coredns/coredns/master.svg)](https://codecov.io/github/coredns/coredns?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/coredns/coredns)](https://goreportcard.com/report/coredns/coredns)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bhttps%3A%2F%2Fgithub.com%2Fcoredns%2Fcoredns.svg?type=shield)](https://app.fossa.io/projects/git%2Bhttps%3A%2F%2Fgithub.com%2Fcoredns%2Fcoredns?ref=badge_shield)

CoreDNS is a DNS server that started as a fork of [Caddy](https://github.com/mholt/caddy/). It has
the same model: it chains middleware. In fact it's so similar that CoreDNS is now a server type
plugin for Caddy.

CoreDNS is also a [Cloud Native Computing Foundation](https://cncf.io) inception level project.

CoreDNS is the successor to [SkyDNS](https://github.com/skynetservices/skydns). SkyDNS is a thin
layer that exposes services in etcd in the DNS. CoreDNS builds on this idea and is a generic DNS
server that can talk to multiple backends (etcd, kubernetes, etc.).

CoreDNS aims to be a fast and flexible DNS server. The keyword here is *flexible*: with CoreDNS you
are able to do what you want with your DNS data. And if not: write some middleware!

CoreDNS can listen for DNS request coming in over UDP/TCP (go'old DNS), TLS ([RFC
7858](https://tools.ietf.org/html/rfc7858)) and gRPC (not a standard).

Currently CoreDNS is able to:

* Serve zone data from a file; both DNSSEC (NSEC only) and DNS are supported (*file*).
* Retrieve zone data from primaries, i.e., act as a secondary server (AXFR only) (*secondary*).
* Sign zone data on-the-fly (*dnssec*).
* Load balancing of responses (*loadbalance*).
* Allow for zone transfers, i.e., act as a primary server (*file*).
* Automatically load zone files from disk (*auto*).
* Caching (*cache*).
* Health checking endpoint (*health*).
* Use etcd as a backend, i.e., a 101.5% replacement for
  [SkyDNS](https://github.com/skynetservices/skydns) (*etcd*).
* Use k8s (kubernetes) as a backend (*kubernetes*).
* Serve as a proxy to forward queries to some other (recursive) nameserver (*proxy*).
* Provide metrics (by using Prometheus) (*metrics*).
* Provide query (*log*) and error (*error*) logging.
* Support the CH class: `version.bind` and friends (*chaos*).
* Profiling support (*pprof*).
* Rewrite queries (qtype, qclass and qname) (*rewrite*).
* Echo back the IP address, transport and port number used (*whoami*).

Each of the middlewares has a README.md of its own.

## Status

CoreDNS can be used as an authoritative nameserver for your domains, and should be stable enough to
provide you with good DNS(SEC) service.

There are still a few known [issues](https://github.com/coredns/coredns/issues), and work is ongoing
on making things fast and to reduce the memory usage.

All in all, CoreDNS should be able to provide you with enough functionality to replace parts of BIND
9, Knot, NSD or PowerDNS and SkyDNS. Most documentation is in the source and some blog articles can
be [found here](https://blog.coredns.io). If you do want to use CoreDNS in production, please
let us know and how we can help.

<https://caddyserver.com/> is also full of examples on how to structure a Corefile (renamed from
Caddyfile when forked).

## Compilation

CoreDNS (as a servertype plugin for Caddy) has a dependency on Caddy, but this is not different than
any other Go dependency. If you have the source of CoreDNS, get all dependencies:

    go get ./...

And then `go build` as you would normally do:

    go build

This should yield a `coredns` binary.

## Examples

When starting CoreDNS without any configuration, it loads the `whoami` middleware and starts
listening on port 53 (override with `-dns.port`), it should show the following:

~~~ txt
.:53
2016/09/18 09:20:50 [INFO] CoreDNS-001
CoreDNS-001
~~~

Any query send to port 53 should return some information; your sending address, port and protocol
used.

If you have a Corefile without a port number specified it will, by default, use port 53, but you
can override the port with the `-dns.port` flag:

`./coredns -dns.port 1053`, runs the server on port 1053.

Start a simple proxy, you'll need to be root to start listening on port 53.

`Corefile` contains:

~~~ txt
.:53 {
    proxy . 8.8.8.8:53
    log stdout
}
~~~

Just start CoreDNS: `./coredns`.
And then just query on that port (53). The query should be forwarded to 8.8.8.8 and the response
will be returned. Each query should also show up in the log.

Serve the (NSEC) DNSSEC-signed `example.org` on port 1053, with errors and logging sent to stdout.
Allow zone transfers to everybody, but specifically mention 1 IP address so that CoreDNS can send
notifies to it.

~~~ txt
example.org:1053 {
    file /var/lib/coredns/example.org.signed {
        transfer to *
        transfer to 2001:500:8f::53
    }
    errors stdout
    log stdout
}
~~~

Serve `example.org` on port 1053, but forward everything that does *not* match `example.org` to a recursive
nameserver *and* rewrite ANY queries to HINFO.

~~~ txt
.:1053 {
    rewrite ANY HINFO
    proxy . 8.8.8.8:53

    file /var/lib/coredns/example.org.signed example.org {
        transfer to *
        transfer to 2001:500:8f::53
    }
    errors stdout
    log stdout
}
~~~

### Zone Specification

The following Corefile fragment is legal, but does not explicitly define a zone to listen on:

~~~ txt
{
   # ...
}
~~~

This defaults to `.:53` (or whatever `-dns.port` is).

The next one only defines a port:
~~~ txt
:123 {
    # ...
}
~~~
This defaults to the root zone `.`, but can't be overruled with the `-dns.port` flag.

Just specifying a zone, default to listening on port 53 (can still be overridden with `-dns.port`):

~~~ txt
example.org {
    # ...
}
~~~

IP addresses are also allowed. They are automatically converted to reverse zones:

~~~ txt
10.0.0.0/24 {
    # ...
}
~~~
Means you are authoritative for `0.0.10.in-addr.arpa.`. 

The netmask must be dividable by 8, if it is not the reverse conversion is not done. This also works
for IPv6 addresses. If for some reason you want to serve a zone named `10.0.0.0/24` add the closing
dot: `10.0.0.0/24.` as this also stops the conversion.

Listening on TLS and for gRPC? Use:

~~~ txt
tls://example.org grpc://example.org {
    # ...
}
~~~

Specifying ports works in the same way:

~~~ txt
grpc://example.org:1443 {
    # ...
}
~~~

When no transport protocol is specified the default `dns://` is assumed.

## Community

- Website: <https://coredns.io>
- Blog: <https://blog.coredns.io>
- Twitter: [@corednsio](https://twitter.com/corednsio)
- Github: <https://github.com/coredns/coredns>
- Mailing list/group: <coredns-discuss@googlegroups.com>
- Slack: #coredns on <https://slack.cncf.io>

## Deployment

Examples for deployment via systemd and other use cases can be found in the
[deployment repository](https://github.com/coredns/deployment).
