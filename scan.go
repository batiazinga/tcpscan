package tcpscan

import (
	"fmt"
	"net"
	"sort"
)

// dialFunc returns true if the port of the host is opened.
type dialFunc func(host string, port int) bool

// worker reads ports to test from the ports channel and send open ports to the collector channel.
// When done, i.e. after the ports channel is closed, it sends a signal to the tracker channel.
func worker(dial dialFunc, host string, ports <-chan int, collector chan<- int, tracker chan<- struct{}) {
	for p := range ports {
		if dial(host, p) {
			collector <- p
		}
	}
	// signal that all ports have processed
	tracker <- struct{}{}
}

func scan(dial dialFunc, host string, numPorts, numWorkers int) []int {
	tracker := make(chan struct{})
	ports := make(chan int, numWorkers)
	collector := make(chan int)
	var openports []int

	// launch all workers in separate goroutines
	for i := 0; i < cap(ports); i++ {
		go worker(dial, host, ports, collector, tracker)
	}

	// start a goroutine which collects open ports
	go func() {
		for port := range collector {
			openports = append(openports, port)
		}
		// signal that all results have been collected
		tracker <- struct{}{}
	}()

	// send ports to dial
	for i := 1; i <= numPorts; i++ {
		ports <- i
	}
	close(ports) // all ports have been sent

	// wait for all workers to be done
	for i := 0; i < numWorkers; i++ {
		<-tracker
	}
	// and since all workers are done, we can close the collector channel
	close(collector)

	// wait for the all results to be collected
	<-tracker

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
