package main

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	if len(os.Args) == 1 {
		fmt.Println("Nothing to ping. Try `goping-cli www.google.com`")
		os.Exit(1)
	}

	start := time.Now()

	var wg sync.WaitGroup
	wg.Add(runtime.NumCPU())

	url := os.Args[1]

	for i := 0; i < runtime.NumCPU(); i++ {

		p := Ping{
			index: i,
			url:   url,
		}

		go func(pinger Ping) {
			defer wg.Done()

			pinger.DoPing()
		}(p)
	}

	// Block until all pings are finished
	wg.Wait()

	fmt.Printf("goping finished in %v\n", time.Since(start))
	os.Exit(0)
}
