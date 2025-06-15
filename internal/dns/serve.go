package dns

import (
	"context"
	"net"
	"strings"
	"time"

	"github.com/ei-sugimoto/wakemae/internal/registry"
	"github.com/miekg/dns"
)

func Serve(addr string, rg *registry.Registry, upstream string) error {
	resolver := &server{rg: rg, upstream: upstream}

	mux := dns.NewServeMux()
	mux.HandleFunc(".", resolver.handle) // catch‑all

	// Start UDP and TCP servers (same mux)
	udpSrv := &dns.Server{Addr: addr, Net: "udp", Handler: mux}
	tcpSrv := &dns.Server{Addr: addr, Net: "tcp", Handler: mux}

	errCh := make(chan error, 2)
	go func() { errCh <- udpSrv.ListenAndServe() }()
	go func() { errCh <- tcpSrv.ListenAndServe() }()

	return <-errCh // first error stops both
}

type server struct {
	rg       *registry.Registry
	upstream string // host:port
}

func (s *server) handle(w dns.ResponseWriter, req *dns.Msg) {
	if len(req.Question) == 0 {
		_ = w.WriteMsg(new(dns.Msg))
		return
	}

	q := req.Question[0]
	fqdn := strings.TrimSuffix(q.Name, ".") // no trailing dot

	// Only handle A / CNAME for PoC
	if ips, ok := s.rg.Resolve(fqdn); ok && q.Qtype == dns.TypeA {
		resp := new(dns.Msg)
		resp.SetReply(req)
		for _, ip := range ips {
			rr := &dns.A{Hdr: dns.RR_Header{Name: q.Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 0}, A: net.ParseIP(ip)}
			resp.Answer = append(resp.Answer, rr)
		}
		_ = w.WriteMsg(resp)
		return
	}

	// Fallback to upstream resolver
	c := dns.Client{Timeout: 2 * time.Second}
	u := s.upstream
	if !strings.Contains(u, ":") { // default port
		u += ":53"
	}
	resp, _, err := c.ExchangeContext(context.Background(), req, u)
	if err == nil {
		_ = w.WriteMsg(resp)
	}
}
