package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

func main() {
	wg := sync.WaitGroup{}
	portsChan := make(chan string)
	resChan := make(chan string)
	argHost := flag.String("h", "", "Host to scan")
	argPortStart := flag.Int("s", 1, "Start port (Default 1)")
	argPortEnd := flag.Int("e", 100, "End port (Default 100)")
	argThreadsNum := flag.Int("t", 4, "Threads count (Default 4)")

	flag.Parse()

	if *argHost == "" {
		flag.Usage()
		os.Exit(1)
	}

	// Pushing events into the channel
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := *argPortStart; i <= *argPortEnd; i++ {
			targetHost := fmt.Sprintf("%s:%d", *argHost, i)
			portsChan <- targetHost
		}
		close(portsChan)
	}()

	// Consuming events
	for i := 0; i <= *argThreadsNum; i++ {
		wg.Add(1)
		go portScanner(portsChan, resChan, &wg)
	}

	// Returning the results
	go func() {
		for result := range resChan {
			fmt.Println(result)
		}
	}()

	wg.Wait()
	close(resChan)
}

// Scan a single port
func portScanner(portsChan, resChan chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for targetHost := range portsChan {
		_, err := net.DialTimeout("tcp", targetHost, time.Millisecond*300)
		if err == nil {
			resChan <- targetHost
		}
	}
}
