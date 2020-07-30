package tcpscan

import (
	"fmt"
	"testing"
)

var (
	allClosed = func(host string, port int) bool { return false }
	allOpen   = func(host string, port int) bool { return true }
)

// TestScan_closed scans closed ports and checks that all ports are detected as closed.
func TestScan_closed(t *testing.T) {
	numPorts := []int{1023, 49151, 65535}
	numWorkers := []int{1, 2, 4, 8, 16, 32, 64}

	for _, p := range numPorts {
		for _, w := range numWorkers {
			t.Run(
				fmt.Sprintf("%d-%d", p, w),
				func(t *testing.T) {
					openPorts := scan(allClosed, "localhost", p, w)

					if len(openPorts) != 0 {
						t.Errorf("%d ports are opened", len(openPorts))
					}
				},
			)
		}
	}
}

// TestScan_open scans open ports and checks that all ports are open.
func TestScan_open(t *testing.T) {
	numPorts := []int{1023, 49151, 65535}
	numWorkers := []int{1, 2, 4, 8, 16, 32, 64}

	for _, p := range numPorts {
		for _, w := range numWorkers {
			t.Run(
				fmt.Sprintf("%d-%d", p, w),
				func(t *testing.T) {
					openPorts := scan(allOpen, "localhost", p, w)

					if len(openPorts) != p {
						t.Errorf("%d/%d ports are opened", len(openPorts), p)
					}
				},
			)
		}
	}
}

// BenchmarkScan_closed benchmarks the scanner with a fast dialer which always finds closed ports.
// This allows evaluating the performance of the concurrency machinery of the scanner without the dial overhead.
func BenchmarkScan_closed(b *testing.B) {
	numPorts := []int{1023, 49151, 65535}
	numWorkers := []int{1, 2, 4, 8, 16, 32, 64}

	for _, p := range numPorts {
		for _, w := range numWorkers {
			b.Run(
				fmt.Sprintf("%d-%d", p, w),
				func(b *testing.B) {
					for i := 0; i < b.N; i++ {
						scan(allClosed, "localhost", p, w)
					}
				},
			)
		}
	}
}

// BenchmarkScan_open benchmarks the scanner with a fast dialer which always finds open ports.
// This is similar to BenchmarkScan_closed but the performance is now impacted by collecting open ports and sorting them.
func BenchmarkScan_open(b *testing.B) {
	numPorts := []int{1023, 49151, 65535}
	numWorkers := []int{1, 2, 4, 8, 16, 32, 64}

	for _, p := range numPorts {
		for _, w := range numWorkers {
			b.Run(
				fmt.Sprintf("%d-%d", p, w),
				func(b *testing.B) {
					for i := 0; i < b.N; i++ {
						scan(allOpen, "localhost", p, w)
					}
				},
			)
		}
	}
}
