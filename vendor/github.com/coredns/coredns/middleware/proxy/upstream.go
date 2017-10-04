package proxy

import (
	"fmt"
	"net"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/coredns/coredns/middleware"
	"github.com/coredns/coredns/middleware/pkg/dnsutil"
	"github.com/coredns/coredns/middleware/pkg/healthcheck"
	"github.com/coredns/coredns/middleware/pkg/tls"
	"github.com/mholt/caddy/caddyfile"
	"github.com/miekg/dns"
)

type staticUpstream struct {
	from string

	healthcheck.HealthCheck

	WithoutPathPrefix string
	IgnoredSubDomains []string
	ex                Exchanger
}

// NewStaticUpstreams parses the configuration input and sets up
// static upstreams for the proxy middleware.
func NewStaticUpstreams(c *caddyfile.Dispenser) ([]Upstream, error) {
	var upstreams []Upstream
	for c.Next() {
		upstream := &staticUpstream{
			from: ".",
			HealthCheck: healthcheck.HealthCheck{
				FailTimeout: 10 * time.Second,
				MaxFails:    1,
				Future:      60 * time.Second,
			},
			ex: newDNSEx(),
		}

		if !c.Args(&upstream.from) {
			return upstreams, c.ArgErr()
		}
		upstream.from = middleware.Host(upstream.from).Normalize()

		to := c.RemainingArgs()
		if len(to) == 0 {
			return upstreams, c.ArgErr()
		}

		// process the host list, substituting in any nameservers in files
		toHosts, err := dnsutil.ParseHostPortOrFile(to...)
		if err != nil {
			return upstreams, err
		}

		for c.NextBlock() {
			if err := parseBlock(c, upstream); err != nil {
				return upstreams, err
			}
		}

		upstream.Hosts = make([]*healthcheck.UpstreamHost, len(toHosts))
		for i, host := range toHosts {
			uh := &healthcheck.UpstreamHost{
				Name:        host,
				Conns:       0,
				Fails:       0,
				FailTimeout: upstream.FailTimeout,

				CheckDown: func(upstream *staticUpstream) healthcheck.UpstreamHostDownFunc {
					return func(uh *healthcheck.UpstreamHost) bool {

						down := false

						uh.CheckMu.Lock()
						until := uh.OkUntil
						uh.CheckMu.Unlock()

						if !until.IsZero() && time.Now().After(until) {
							down = true
						}

						fails := atomic.LoadInt32(&uh.Fails)
						if fails >= upstream.MaxFails && upstream.MaxFails != 0 {
							down = true
						}
						return down
					}
				}(upstream),
				WithoutPathPrefix: upstream.WithoutPathPrefix,
			}

			upstream.Hosts[i] = uh
		}
		upstream.Start()

		upstreams = append(upstreams, upstream)
	}
	return upstreams, nil
}

func (u *staticUpstream) From() string {
	return u.from
}

func parseBlock(c *caddyfile.Dispenser, u *staticUpstream) error {
	switch c.Val() {
	case "policy":
		if !c.NextArg() {
			return c.ArgErr()
		}
		policyCreateFunc, ok := healthcheck.SupportedPolicies[c.Val()]
		if !ok {
			return c.ArgErr()
		}
		u.Policy = policyCreateFunc()
	case "fail_timeout":
		if !c.NextArg() {
			return c.ArgErr()
		}
		dur, err := time.ParseDuration(c.Val())
		if err != nil {
			return err
		}
		u.FailTimeout = dur
	case "max_fails":
		if !c.NextArg() {
			return c.ArgErr()
		}
		n, err := strconv.Atoi(c.Val())
		if err != nil {
			return err
		}
		u.MaxFails = int32(n)
	case "health_check":
		if !c.NextArg() {
			return c.ArgErr()
		}
		var err error
		u.HealthCheck.Path, u.HealthCheck.Port, err = net.SplitHostPort(c.Val())
		if err != nil {
			return err
		}
		u.HealthCheck.Interval = 30 * time.Second
		if c.NextArg() {
			dur, err := time.ParseDuration(c.Val())
			if err != nil {
				return err
			}
			u.HealthCheck.Interval = dur
			u.Future = 2 * dur

			// set a minimum of 3 seconds
			if u.Future < (3 * time.Second) {
				u.Future = 3 * time.Second
			}
		}
	case "without":
		if !c.NextArg() {
			return c.ArgErr()
		}
		u.WithoutPathPrefix = c.Val()
	case "except":
		ignoredDomains := c.RemainingArgs()
		if len(ignoredDomains) == 0 {
			return c.ArgErr()
		}
		for i := 0; i < len(ignoredDomains); i++ {
			ignoredDomains[i] = middleware.Host(ignoredDomains[i]).Normalize()
		}
		u.IgnoredSubDomains = ignoredDomains
	case "spray":
		u.Spray = &healthcheck.Spray{}
	case "protocol":
		encArgs := c.RemainingArgs()
		if len(encArgs) == 0 {
			return c.ArgErr()
		}
		switch encArgs[0] {
		case "dns":
			if len(encArgs) > 1 {
				if encArgs[1] == "force_tcp" {
					opts := Options{ForceTCP: true}
					u.ex = newDNSExWithOption(opts)
				} else {
					return fmt.Errorf("only force_tcp allowed as parameter to dns")
				}
			} else {
				u.ex = newDNSEx()
			}
		case "https_google":
			boot := []string{"8.8.8.8:53", "8.8.4.4:53"}
			if len(encArgs) > 2 && encArgs[1] == "bootstrap" {
				boot = encArgs[2:]
			}

			u.ex = newGoogle("", boot) // "" for default in google.go
		case "grpc":
			if len(encArgs) == 2 && encArgs[1] == "insecure" {
				u.ex = newGrpcClient(nil, u)
				return nil
			}
			tls, err := tls.NewTLSConfigFromArgs(encArgs[1:]...)
			if err != nil {
				return err
			}
			u.ex = newGrpcClient(tls, u)
		default:
			return fmt.Errorf("%s: %s", errInvalidProtocol, encArgs[0])
		}

	default:
		return c.Errf("unknown property '%s'", c.Val())
	}
	return nil
}

func (u *staticUpstream) IsAllowedDomain(name string) bool {
	if dns.Name(name) == dns.Name(u.From()) {
		return true
	}

	for _, ignoredSubDomain := range u.IgnoredSubDomains {
		if middleware.Name(ignoredSubDomain).Matches(name) {
			return false
		}
	}
	return true
}

func (u *staticUpstream) Exchanger() Exchanger { return u.ex }
