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

func get8uint32(in []byte, out []uint32, ctrl uint16) int {
	// bounds check to prevent undefined behavior
	_ = in[7] // best effort, there should be at least 8 bytes
	_ = out[7]
	return int(get8uint32Fast(in, out, ctrl, shared.DecodeShuffleTable))
}

func get8uint32Fast(in []byte, out []uint32, ctrl uint16, shuffle *[256][16]uint8) uint64