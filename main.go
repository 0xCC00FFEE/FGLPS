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
	resultsChan := make(chan string)

	argHost := flag.String("host", "", "Host to scan")
	argFirstPort := flag.Int("firstPort", 1, "First port of port range to scan (1-65535)")
	argLastPort := flag.Int("lastPort", 65535, "Last port of port range to scan (1-65535)")
	argThreadsNum := flag.Int("threads", 65535, "Thread count. (maximum simultaneous port scans)")
	argPortTimeout := flag.Int("portTimeout", 5, "Port timeout in seconds.")
	flag.Parse()

	showUsageInfo := false
	if *argHost == "" {
		showUsageInfo = true
	}
	if *argFirstPort > *argLastPort {
		fmt.Println("ERROR: -firstPort cannot be > -lastPort")
		showUsageInfo = true
	}
	if *argFirstPort < 1 || *argFirstPort > 65535 {
		fmt.Println("ERROR: -firstPort must be >= 1 and <= 65535")
		showUsageInfo = true
	}
	if *argLastPort < 1 || *argLastPort > 65535 {
		fmt.Println("ERROR: -lastPort must be >= 1 and <= 65535")
		showUsageInfo = true
	}
	if *argThreadsNum < 1 || *argThreadsNum > 65535 {
		fmt.Println("ERROR: -threads must be >= 1 and <= 65535")
		showUsageInfo = true
	}
	if *argPortTimeout < 1 || *argPortTimeout > 65535 {
		fmt.Println("ERROR: -portTimeout must be >= 1 and <= 65535")
		showUsageInfo = true
	}

	if showUsageInfo {
		showUsage()
	}

	// Pushing events into the channel
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := *argFirstPort; i <= *argLastPort; i++ {
			targetHost := fmt.Sprintf("%s:%d", *argHost, i)
			portsChan <- targetHost
		}
		close(portsChan)
	}()

	// Consuming events
	portRangeCount := *argLastPort - *argFirstPort + 1
	maxParallelPortScans := *argThreadsNum
	if portRangeCount < *argThreadsNum {
		maxParallelPortScans = portRangeCount
	}
	for i := 1; i <= maxParallelPortScans; i++ {
		wg.Add(1)
		go portScanner(portsChan, resultsChan, &wg, *argPortTimeout)
	}

	// Returning the results
	go func() {
		for result := range resultsChan {
			fmt.Println(result)
		}
	}()

	wg.Wait()
	close(resultsChan)
}

func showUsage() {
	fmt.Println("")
	fmt.Println("Purpose: Rapidly scans a TCP port range for open ports.")
	fmt.Println("")
	flag.Usage()
	os.Exit(1)
}

// Scan a single port
func portScanner(portsChan, resultsChan chan string, wg *sync.WaitGroup, portTimeout int) {
	defer wg.Done()
	for targetHost := range portsChan {
		_, err := net.DialTimeout("tcp", targetHost, time.Second*time.Duration(portTimeout))
		if err == nil {
			resultsChan <- targetHost
		}
	}
}
