package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/miekg/dns"
)

var recursor string
var ctf_domain string

func parseQuery(m *dns.Msg, ip string) {
	for _, q := range m.Question {
		switch q.Qtype {
		case dns.TypeA:
			log.Printf("Query for %s\n", q.Name)
			s := strings.Split(ip, ".")
			i, _ := strconv.Atoi(s[3])
			s[3] = strconv.Itoa(i - 1)
			rr, err := dns.NewRR(fmt.Sprintf("%s A %s", q.Name, strings.Join(s, ".")))
			if err == nil {
				m.Answer = append(m.Answer, rr)
			}
		}
	}
}

func handleCtfDnsRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	switch r.Opcode {
	case dns.OpcodeQuery:
		parseQuery(m, strings.Split(w.RemoteAddr().String(), ":")[0])
	}

	w.WriteMsg(m)
}

func handleOtherDnsRequest(resp dns.ResponseWriter, req *dns.Msg) {
	if len(req.Question) == 0 {
		respond(resp, req, dns.RcodeFormatError)
		return
	}

	network := "udp"
	if _, ok := resp.RemoteAddr().(*net.TCPAddr); ok {
		network = "tcp"
	}

	c := &dns.Client{Net: network}
	r, _, err := c.Exchange(req, recursor)
	if err == nil {
		log.Printf("[info] using %s to answer %s", recursor, req.Question[0].Name)
		if err := resp.WriteMsg(r); err != nil {
			log.Printf("[WARN] dns: failed to respond: %v", err)
		}
		return
	}

	respond(resp, req, dns.RcodeServerFailure)
}

func respond(resp dns.ResponseWriter, req *dns.Msg, code int) {
	m := &dns.Msg{}
	m.SetReply(req)
	m.RecursionAvailable = true
	m.SetRcode(req, code)
	resp.WriteMsg(m)
}

func main() {
	var ctf_domain string
	flag.StringVar(&ctf_domain, "ctf_domain", "", "CTF Domain")
	flag.StringVar(&recursor, "dns_server", "8.8.8.8:53", "External DNS Resolver")
	flag.Parse()
	if ctf_domain == "" {
		panic("Please set ctf_domain")
	}
	dns.HandleFunc(ctf_domain+".", handleCtfDnsRequest)
	dns.HandleFunc(".", handleOtherDnsRequest)

	port := 53
	server := &dns.Server{Addr: ":" + strconv.Itoa(port), Net: "udp"}
	log.Printf("Starting at %d\n", port)
	err := server.ListenAndServe()
	defer server.Shutdown()
	if err != nil {
		log.Fatalf("Failed to start server: %s\n ", err.Error())
	}
}
