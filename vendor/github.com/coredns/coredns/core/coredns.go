// Package core registers the server and all plugins we support.
package core

import (
	// plug in the server
	_ "github.com/coredns/coredns/core/dnsserver"

	// plug in the standard directives (sorted)
	_ "github.com/coredns/coredns/middleware/auto"
	_ "github.com/coredns/coredns/middleware/bind"
	_ "github.com/coredns/coredns/middleware/cache"
	_ "github.com/coredns/coredns/middleware/chaos"
	_ "github.com/coredns/coredns/middleware/dnssec"
	_ "github.com/coredns/coredns/middleware/dnstap"
	_ "github.com/coredns/coredns/middleware/erratic"
	_ "github.com/coredns/coredns/middleware/errors"
	_ "github.com/coredns/coredns/middleware/etcd"
	_ "github.com/coredns/coredns/middleware/file"
	_ "github.com/coredns/coredns/middleware/health"
	_ "github.com/coredns/coredns/middleware/kubernetes"
	_ "github.com/coredns/coredns/middleware/loadbalance"
	_ "github.com/coredns/coredns/middleware/log"
	_ "github.com/coredns/coredns/middleware/metrics"
	_ "github.com/coredns/coredns/middleware/pprof"
	_ "github.com/coredns/coredns/middleware/proxy"
	_ "github.com/coredns/coredns/middleware/reverse"
	_ "github.com/coredns/coredns/middleware/rewrite"
	_ "github.com/coredns/coredns/middleware/root"
	_ "github.com/coredns/coredns/middleware/secondary"
	_ "github.com/coredns/coredns/middleware/trace"
	_ "github.com/coredns/coredns/middleware/whoami"
)
