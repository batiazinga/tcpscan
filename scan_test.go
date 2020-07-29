package tcpscan_test

import (
	"fmt"
	"testing"

	"github.com/batiazinga/tcpscan"
)

func BenchmarkScan(b *testing.B) {
	host := "127.0.0.1"
	numPorts := []int{1023, 49151, 65535}
	numWorkers := []int{1, 2, 4, 8, 16, 32, 64}

	for _, p := range numPorts {
		for _, w := range numWorkers {
			b.Run(
				fmt.Sprintf("%d-%d", p, w),
				func(b *testing.B) {
					for i := 0; i < b.N; i++ {
						tcpscan.Scan(host, p, w)
					}
				},
			)
		}
	}
}
