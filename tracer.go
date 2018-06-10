package main

import (
	"crypto/tls"
	"io"
	"net/http/httptrace"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	jaeger "github.com/uber/jaeger-client-go"
	jaegerlog "github.com/uber/jaeger-client-go/log"
)

// InitGlobalTracer configures an open tracing tracer
func InitGlobalTracer() (opentracing.Tracer, io.Closer) {
	// Configure jaeger tracer
	sender, _ := jaeger.NewUDPTransport("localhost:6831", 0)
	sampler := jaeger.NewConstSampler(true)
	options := jaeger.ReporterOptions.BufferFlushInterval(1 * time.Second)
	reporter := jaeger.NewRemoteReporter(sender, options)
	logger := jaeger.NewLoggingReporter(jaegerlog.StdLogger)

	tracer, closer := jaeger.NewTracer(
		"goping",
		sampler,
		jaeger.NewCompositeReporter(logger, reporter),
	)

	opentracing.SetGlobalTracer(tracer)

	return tracer, closer
}

// NewClientTrace http request trace
func NewClientTrace(span opentracing.Span) *httptrace.ClientTrace {
	trace := &clientTrace{span: span}
	return &httptrace.ClientTrace{
		DNSStart:             trace.dnsStart,
		DNSDone:              trace.dnsDone,
		GetConn:              trace.getConn,
		GotConn:              trace.gotConn,
		GotFirstResponseByte: trace.gotFirstResponseByte,
		TLSHandshakeDone:     trace.tlsHandshakeDone,
	}
}

type clientTrace struct {
	span opentracing.Span
}

func (h *clientTrace) dnsStart(info httptrace.DNSStartInfo) {
	h.span.LogFields(
		log.String("event", "DNSStart"),
		log.Object("host", info.Host),
	)
}

func (h *clientTrace) dnsDone(info httptrace.DNSDoneInfo) {
	h.span.LogFields(
		log.String("event", "DNSDone"),
		log.Object("addrs", info.Addrs),
		log.Bool("coalesced", info.Coalesced),
		log.Error(info.Err),
	)
}

func (h *clientTrace) getConn(hostPort string) {
	h.span.LogFields(
		log.String("event", "GetConn"),
		log.String("port", hostPort),
	)
}

func (h *clientTrace) gotConn(info httptrace.GotConnInfo) {
	h.span.LogFields(
		log.String("event", "GotConn"),
		log.Object("time", info.IdleTime),
		log.Object("con", info.Conn),
		log.Object("idle.time", info.IdleTime),
		log.Bool("was.idle", info.WasIdle),
	)
}

func (h *clientTrace) gotFirstResponseByte() {
	h.span.LogFields(log.String("event", "GotFirstResponseByte"))
}

func (h *clientTrace) tlsHandshakeDone(state tls.ConnectionState, e error) {
	h.span.LogFields(
		log.String("event", "TLSHandshakeDone"),
		log.Object("state", state),
		log.Error(e),
	)
}
