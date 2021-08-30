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
	return int(get8uint32Fast(in, out, ctrl,
		shared.DecodeShuffleTable,
		shared.PerControlLenTable,
	))
}

func get8uint32Fast(
	in []byte, out []uint32, ctrl uint16,
	shuffle *[256][16]uint8, lenTable *[256]uint8,
) (r uint64)
