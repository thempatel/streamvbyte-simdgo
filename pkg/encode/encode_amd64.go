// +build amd64

package encode

import (
	"github.com/theMPatel/streamvbyte-simdgo/pkg/shared"
	"golang.org/x/sys/cpu"
)

func GetMode() shared.PerformanceMode {
	if cpu.X86.HasAVX && cpu.X86.HasAVX2 {
		return shared.Fast
	}
	return shared.Normal
}

func put8uint32(in []uint32, out []byte) uint16 {
	return put8uint32Fast(in, out,
		shared.EncodeShuffleTable,
		shared.PerControlLenTable,
	)
}

func put8uint32Fast(
	in []uint32, outBytes []byte,
	shuffle *[256][16]uint8, lenTable *[256]uint8,
) (r uint16)
