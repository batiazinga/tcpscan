package tcpscan

import (
	"fmt"
	"testing"
)

// TestScan_closed scans closed ports and checks that all ports are detected as closed.
func TestScan_closed(t *testing.T) {
	allClosed := func(host string, port int) bool { return false }

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
	allOpen := func(host string, port int) bool { return true }

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
