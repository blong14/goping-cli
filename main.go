package main

import (
	"fmt"
	"os"
	"runtime"
	"sync"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

var (
	parentCtx opentracing.SpanContext
	wg        sync.WaitGroup
)

func main() {
	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus)

	if len(os.Args) == 1 {
		fmt.Println("Nothing to ping. Try `goping-cli www.google.com`")
		os.Exit(1)
	}

	wg.Add(runtime.NumCPU())

	url := os.Args[1]

	tracer, closer := InitGlobalTracer()

	parentSpan := tracer.StartSpan(
		url,
		opentracing.Tag{Key: "runtime.num.cpu", Value: cpus},
		opentracing.Tag{Key: "runtime.goos", Value: runtime.GOOS},
		opentracing.Tag{Key: "runtime.goarch", Value: runtime.GOARCH},
		opentracing.Tag{Key: "runtime.version", Value: runtime.Version()},
	)

	for i := 0; i < cpus; i++ {

		p := Ping{index: i, url: url}

		go func(pinger Ping) {
			defer wg.Done()

			parentSpan.LogFields(
				log.String("event", "start ping"),
				log.Int("ping.index", pinger.index),
			)

			pinger.DoPing(parentSpan.Context())
		}(p)
	}

	// Block until all pings are finished
	wg.Wait()

	parentSpan.Finish()

	closer.Close()

	os.Exit(0)
}
