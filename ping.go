package main

import (
	"net/http"
	"net/http/httptrace"
	"regexp"

	opentracing "github.com/opentracing/opentracing-go"
)

// Ping data associated with a single ping
type Ping struct {
	index int
	url   string
	err   error
}

// DoPing pings the given url
func (ping *Ping) DoPing(ctx opentracing.SpanContext) {

	match, _ := regexp.MatchString("htt([ps]+)://", ping.url)

	if !match {
		ping.url = "http://" + ping.url
	}

	tracer := opentracing.GlobalTracer()

	// start a new Span to wrap HTTP request
	span := tracer.StartSpan(
		"DoPing",
		opentracing.ChildOf(ctx),
		opentracing.Tag{Key: "url", Value: ping.url},
		opentracing.Tag{Key: "ping.index", Value: ping.index},
	)

	// make sure the Span is finished once we're done
	defer span.Finish()

	req, _ := http.NewRequest("GET", ping.url, nil)

	trace := NewClientTrace(span)

	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))

	_, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		ping.err = err
	}

	span.SetTag("error", ping.err)
}
