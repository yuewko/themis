package dnssec

import (
	"strings"
	"testing"

	"github.com/coredns/coredns/middleware/file"
	"github.com/coredns/coredns/middleware/pkg/cache"
	"github.com/coredns/coredns/middleware/pkg/dnsrecorder"
	"github.com/coredns/coredns/middleware/test"

	"github.com/miekg/dns"
	"golang.org/x/net/context"
)

var dnssecTestCases = []test.Case{
	{
		Qname: "miek.nl.", Qtype: dns.TypeDNSKEY,
		Answer: []dns.RR{
			test.DNSKEY("miek.nl.	3600	IN	DNSKEY	257 3 13 0J8u0XJ9GNGFEBXuAmLu04taHG4"),
		},
	},
	{
		Qname: "miek.nl.", Qtype: dns.TypeDNSKEY, Do: true,
		Answer: []dns.RR{
			test.DNSKEY("miek.nl.	3600	IN	DNSKEY	257 3 13 0J8u0XJ9GNGFEBXuAmLu04taHG4"),
			test.RRSIG("miek.nl.	3600	IN	RRSIG	DNSKEY 13 2 3600 20160503150844 20160425120844 18512 miek.nl. Iw/kNOyM"),
		},
		Extra: []dns.RR{test.OPT(4096, true)},
	},
}

var dnsTestCases = []test.Case{
	{
		Qname: "miek.nl.", Qtype: dns.TypeDNSKEY,
		Answer: []dns.RR{
			test.DNSKEY("miek.nl.	3600	IN	DNSKEY	257 3 13 0J8u0XJ9GNGFEBXuAmLu04taHG4"),
		},
	},
	{
		Qname: "miek.nl.", Qtype: dns.TypeMX,
		Answer: []dns.RR{
			test.MX("miek.nl.	1800	IN	MX	1 aspmx.l.google.com."),
		},
		Ns: []dns.RR{
			test.NS("miek.nl.	1800	IN	NS	linode.atoom.net."),
		},
	},
	{
		Qname: "miek.nl.", Qtype: dns.TypeMX, Do: true,
		Answer: []dns.RR{
			test.MX("miek.nl.	1800	IN	MX	1 aspmx.l.google.com."),
			test.RRSIG("miek.nl.	1800	IN	RRSIG	MX 13 2 3600 20160503192428 20160425162428 18512 miek.nl. 4nxuGKitXjPVA9zP1JIUvA09"),
		},
		Ns: []dns.RR{
			test.NS("miek.nl.	1800	IN	NS	linode.atoom.net."),
			test.RRSIG("miek.nl.	1800	IN	RRSIG	NS 13 2 3600 20161217114912 20161209084912 18512 miek.nl. ad9gA8VWgF1H8ze9/0Rk2Q=="),
		},
		Extra: []dns.RR{test.OPT(4096, true)},
	},
	{
		Qname: "www.miek.nl.", Qtype: dns.TypeAAAA, Do: true,
		Answer: []dns.RR{
			test.AAAA("a.miek.nl.	1800	IN	AAAA	2a01:7e00::f03c:91ff:fef1:6735"),
			test.RRSIG("a.miek.nl.	1800	IN	RRSIG	AAAA 13 3 3600 20160503193047 20160425163047 18512 miek.nl. UAyMG+gcnoXW3"),
			test.CNAME("www.miek.nl.	1800	IN	CNAME	a.miek.nl."),
			test.RRSIG("www.miek.nl.	1800	IN	RRSIG	CNAME 13 3 3600 20160503193047 20160425163047 18512 miek.nl. E3qGZn"),
		},
		Ns: []dns.RR{
			test.NS("miek.nl.	1800	IN	NS	linode.atoom.net."),
			test.RRSIG("miek.nl.	1800	IN	RRSIG	NS 13 2 3600 20161217114912 20161209084912 18512 miek.nl. ad9gA8VWgF1H8ze9/0Rk2Q=="),
		},
		Extra: []dns.RR{test.OPT(4096, true)},
	},
	{
		Qname: "www.example.org.", Qtype: dns.TypeAAAA, Do: true,
		Rcode: dns.RcodeServerFailure,
		// Extra: []dns.RR{test.OPT(4096, true)}, // test.ErrorHandler is a simple handler that does not do EDNS.
	},
}

func TestLookupZone(t *testing.T) {
	zone, err := file.Parse(strings.NewReader(dbMiekNL), "miek.nl.", "stdin", 0)
	if err != nil {
		return
	}
	fm := file.File{Next: test.ErrorHandler(), Zones: file.Zones{Z: map[string]*file.Zone{"miek.nl.": zone}, Names: []string{"miek.nl."}}}
	dnskey, rm1, rm2 := newKey(t)
	defer rm1()
	defer rm2()
	c := cache.New(defaultCap)
	dh := New([]string{"miek.nl."}, []*DNSKEY{dnskey}, fm, c)
	ctx := context.TODO()

	for _, tc := range dnsTestCases {
		m := tc.Msg()

		rec := dnsrecorder.New(&test.ResponseWriter{})
		_, err := dh.ServeDNS(ctx, rec, m)
		if err != nil {
			t.Errorf("expected no error, got %v\n", err)
			return
		}

		resp := rec.Msg
		test.SortAndCheck(t, resp, tc)
	}
}

func TestLookupDNSKEY(t *testing.T) {
	dnskey, rm1, rm2 := newKey(t)
	defer rm1()
	defer rm2()
	c := cache.New(defaultCap)
	dh := New([]string{"miek.nl."}, []*DNSKEY{dnskey}, test.ErrorHandler(), c)
	ctx := context.TODO()

	for _, tc := range dnssecTestCases {
		m := tc.Msg()

		rec := dnsrecorder.New(&test.ResponseWriter{})
		_, err := dh.ServeDNS(ctx, rec, m)
		if err != nil {
			t.Errorf("expected no error, got %v\n", err)
			return
		}

		resp := rec.Msg
		if !resp.Authoritative {
			t.Errorf("Authoritative Answer should be true, got false")
		}

		test.SortAndCheck(t, resp, tc)
	}
}

const dbMiekNL = `
$TTL    30M
$ORIGIN miek.nl.
@       IN      SOA     linode.atoom.net. miek.miek.nl. (
                             1282630057 ; Serial
                             4H         ; Refresh
                             1H         ; Retry
                             7D         ; Expire
                             4H )       ; Negative Cache TTL
                IN      NS      linode.atoom.net.

                IN      MX      1  aspmx.l.google.com.

                IN      A       139.162.196.78
                IN      AAAA    2a01:7e00::f03c:91ff:fef1:6735

a               IN      A       139.162.196.78
                IN      AAAA    2a01:7e00::f03c:91ff:fef1:6735
www             IN      CNAME   a`
