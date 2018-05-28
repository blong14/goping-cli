package main

import (
	"fmt"
	"net/http"
	"net/http/httptrace"
	"regexp"
	"time"
)

// Ping data associated with a single ping
type Ping struct {
	firstByte time.Duration
	connStart time.Duration
	connDone  time.Duration
	dnsStart  time.Duration
	dnsDone   time.Duration
	index     int
	url       string
	err       error
}

// DoPing pings the given url
func (ping *Ping) DoPing() {

	match, _ := regexp.MatchString("htt([ps]+)://", ping.url)

	if !match {
		ping.url = "https://" + ping.url
	}

	start := time.Now()

	req, _ := http.NewRequest("GET", ping.url, nil)

	trace := &httptrace.ClientTrace{
		GotFirstResponseByte: func() {
			ping.firstByte = time.Now().Sub(start)
		},
		GetConn: func(hostPort string) {
			ping.connStart = time.Now().Sub(start)
		},
		GotConn: func(connInfo httptrace.GotConnInfo) {
			ping.connDone = time.Now().Sub(start)
		},
		DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
			ping.dnsDone = time.Now().Sub(start)
		},
		DNSStart: func(dnsStartInfo httptrace.DNSStartInfo) {
			ping.dnsStart = time.Now().Sub(start)
		},
	}

	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))

	_, err := http.DefaultTransport.RoundTrip(req)

	if err != nil {
		ping.err = err
	}

	fmt.Println(ping)
}

// String toString
func (ping *Ping) String() string {
	return fmt.Sprintf(
		"URL: %s\nIndex: %d\nConn Start: %v\nConn Done: %v\nDNS Start: %v\nDNS End: %v\nFirst Byte: %v\nErrors: %v\n\n",
		ping.url,
		ping.index,
		ping.connStart,
		ping.connDone,
		ping.dnsStart,
		ping.dnsDone,
		ping.firstByte,
		ping.err,
	)
}
