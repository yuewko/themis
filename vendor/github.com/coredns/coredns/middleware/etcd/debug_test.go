// +build etcd

package etcd

import (
	"testing"

	"github.com/coredns/coredns/middleware/etcd/msg"
	"github.com/coredns/coredns/middleware/pkg/dnsrecorder"
	"github.com/coredns/coredns/middleware/test"

	"github.com/miekg/dns"
)

func TestDebugLookup(t *testing.T) {
	etc := newEtcdMiddleware()
	etc.Debugging = true

	for _, serv := range servicesDebug {
		set(t, etc, serv.Key, 0, serv)
		defer delete(t, etc, serv.Key)
	}

	for _, tc := range dnsTestCasesDebug {
		m := tc.Msg()

		rec := dnsrecorder.New(&test.ResponseWriter{})
		etc.ServeDNS(ctxt, rec, m)

		resp := rec.Msg
		test.SortAndCheck(t, resp, tc)
	}
}

func TestDebugLookupFalse(t *testing.T) {
	etc := newEtcdMiddleware()

	for _, serv := range servicesDebug {
		set(t, etc, serv.Key, 0, serv)
		defer delete(t, etc, serv.Key)
	}
	for _, tc := range dnsTestCasesDebugFalse {
		m := tc.Msg()

		rec := dnsrecorder.New(&test.ResponseWriter{})
		etc.ServeDNS(ctxt, rec, m)

		resp := rec.Msg
		test.SortAndCheck(t, resp, tc)
	}
}

var servicesDebug = []*msg.Service{
	{Host: "127.0.0.1", Key: "a.dom.skydns.test."},
	{Host: "127.0.0.2", Key: "b.sub.dom.skydns.test."},
}

var dnsTestCasesDebug = []test.Case{
	{
		Qname: "o-o.debug.dom.skydns.test.", Qtype: dns.TypeA,
		Answer: []dns.RR{
			test.A("dom.skydns.test. 300 IN A 127.0.0.1"),
			test.A("dom.skydns.test. 300 IN A 127.0.0.2"),
		},
		Extra: []dns.RR{
			test.TXT(`a.dom.skydns.test. 300	CH	TXT	"127.0.0.1:0(10,0,,false)[0,]"`),
			test.TXT(`b.sub.dom.skydns.test. 300	CH	TXT	"127.0.0.2:0(10,0,,false)[0,]"`),
		},
	},
	{
		Qname: "o-o.debug.dom.skydns.test.", Qtype: dns.TypeTXT,
		Ns: []dns.RR{
			test.SOA("skydns.test. 300 IN SOA ns.dns.skydns.test. hostmaster.skydns.test. 1463943291 7200 1800 86400 60"),
		},
		Extra: []dns.RR{
			test.TXT(`a.dom.skydns.test. 300	CH	TXT	"127.0.0.1:0(10,0,,false)[0,]"`),
			test.TXT(`b.sub.dom.skydns.test. 300	CH	TXT	"127.0.0.2:0(10,0,,false)[0,]"`),
		},
	},
}

var dnsTestCasesDebugFalse = []test.Case{
	{
		Qname: "o-o.debug.dom.skydns.test.", Qtype: dns.TypeA,
		Rcode: dns.RcodeNameError,
		Ns: []dns.RR{
			test.SOA("skydns.test. 300 IN SOA ns.dns.skydns.test. hostmaster.skydns.test. 1463943291 7200 1800 86400 60"),
		},
	},
	{
		Qname: "o-o.debug.dom.skydns.test.", Qtype: dns.TypeTXT,
		Rcode: dns.RcodeNameError,
		Ns: []dns.RR{
			test.SOA("skydns.test. 300 IN SOA ns.dns.skydns.test. hostmaster.skydns.test. 1463943291 7200 1800 86400 60"),
		},
	},
}
