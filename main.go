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

	argHost := flag.String("host", "", "Host to scan")
	argFirstPort := flag.Int("firstPort", 1, "First port of port range to scan (1-65535)")
	argLastPort := flag.Int("lastPort", 65535, "Last port of port range to scan (1-65535)")
	argThreadsNum := flag.Int("threads", 65535, "Thread count. (maximum simultaneous port scans)")
	argPortTimeout := flag.Int("portTimeout", 5, "Port timeout in seconds.")
	flag.Parse()

	showUsageInfo := false
	if *argHost == "" {
		showUsageInfo = true
	} else {
		err := checkForInvalidHostname(*argHost)
		if err != nil {
			fmt.Printf("ERROR: 'Hostname Test Failed' ... %v\n", err)
		}
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
		showUsageAndExitAbnormally()
	}

	fmt.Println("# INFO: Starting scan. Open ports will be listed below, in random order (for performance reasons).")

	for i := *argFirstPort; i <= *argLastPort; i++ {
		targetHostWithPort := fmt.Sprintf("%s:%d", *argHost, i)

		// increment the wait group value to hold it open while the goroutine is still running
		wg.Add(1)
		go func() {
			// decrement the wait group value upon completion of the individual port scan
			defer wg.Done()

			if isTcpPortOpen(targetHostWithPort, *argPortTimeout) {
				fmt.Println(targetHostWithPort)
			}
		}()
	}

	wg.Wait()
}

func checkForInvalidHostname(hostname string) error {
	_, err := net.LookupHost(hostname)
	return err
}

func showUsageAndExitAbnormally() {
	fmt.Println("")
	fmt.Println("Purpose: Rapidly scans a TCP port range for open ports.")
	fmt.Println("")
	flag.Usage()
	os.Exit(1)
}

// isTcpPortOpen returns true if the port accepts a TCP connection
func isTcpPortOpen(targetHostWithPort string, portTimeout int) bool {
	_, err := net.DialTimeout("tcp", targetHostWithPort, time.Second*time.Duration(portTimeout))
	if err == nil {
		return true
	}
	return false
}
