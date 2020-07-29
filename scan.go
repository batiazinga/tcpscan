package tcpscan

import (
	"fmt"
	"net"
	"sort"
)

// dialFunc returns true if the port of the host is opened.
type dialFunc func(host string, port int) bool

func worker(dial dialFunc, host string, ports <-chan int, results chan<- int) {
	for p := range ports {
		if dial(host, p) {
			results <- p
		} else {
			results <- 0
		}
	}
}

func scan(dial dialFunc, host string, numPorts, numWorkers int) []int {
	ports := make(chan int, numWorkers)
	results := make(chan int)
	var openports []int

	for i := 0; i < cap(ports); i++ {
		go worker(dial, host, ports, results)
	}

	go func() {
		for i := 1; i <= numPorts; i++ {
			ports <- i
		}
	}()

	for i := 0; i < numPorts; i++ {
		port := <-results
		if port != 0 {
			openports = append(openports, port)
		}
	}

	close(ports)
	close(results)
	sort.Ints(openports)
	return openports
}

// Scan scans ports from 1 to numPorts (included)
// using numWorkers in parallel.
// It returns the sorted list of open ports.
func Scan(host string, numPorts, numWorkers int) []int {
	dial := func(host string, port int) bool {
		address := host + fmt.Sprintf(":%d", port)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			return false
		}
		conn.Close()
		return true
	}

	return scan(dial, host, numPorts, numWorkers)
}
