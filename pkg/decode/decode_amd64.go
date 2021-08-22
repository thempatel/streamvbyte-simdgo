// +build amd64

package decode

import (
	"github.com/theMPatel/streamvbyte-simdgo/pkg/shared"
	"golang.org/x/sys/cpu"
)

func GetMode() shared.PerformanceMode {
	if cpu.X86.HasAVX {
		return shared.Fast
	}
	return shared.Normal
}

func get8uint32(in []byte, out []uint32, ctrl uint8) int {
	return 0
}