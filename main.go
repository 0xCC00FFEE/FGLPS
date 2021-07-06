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

	argHost := flag.String("host", "", "Host to scan")
	argFirstPort := flag.Int("firstPort", 1, "First port of port range to scan (1-65535)")
	argLastPort := flag.Int("lastPort", 100, "Last port of port range to scan (1-65535)")
	argThreadsNum := flag.Int("threads", 65535, "Thread count. (maximum simultaneous port scans)")
	argPortTimeout := flag.Int("portTimeout", 5, "Port timeout in seconds.")
	flag.Parse()

	showUsageInfo := false
	if *argHost == "" {
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
	for i := 0; i <= *argThreadsNum; i++ {
		wg.Add(1)
		go portScanner(portsChan, resChan, &wg, *argPortTimeout)
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

func showUsage() {
	fmt.Println("")
	fmt.Println("Purpose: Rapidly TCP scans a port range for open ports.")
	fmt.Println("")
	flag.Usage()
	os.Exit(1)
}

// Scan a single port
func portScanner(portsChan, resChan chan string, wg *sync.WaitGroup, portTimeout int) {
	defer wg.Done()
	for targetHost := range portsChan {
		_, err := net.DialTimeout("tcp", targetHost, time.Second*time.Duration(portTimeout))
		if err == nil {
			resChan <- targetHost
		}
	}
}
