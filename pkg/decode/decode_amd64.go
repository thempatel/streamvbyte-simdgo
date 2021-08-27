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
	lower, upper := uint8(ctrl&0xff), uint8(ctrl>>8)
	sizeLower, sizeUpper := shared.PerControlLenTable[lower], shared.PerControlLenTable[upper]
	get8uint32Fast(in, out,
		&shared.DecodeShuffleTable[lower],
		&shared.DecodeShuffleTable[upper],
		sizeLower,
	)
	return int(sizeLower+sizeUpper)
}

func get8uint32Fast(in []byte, out []uint32, shufA *[16]uint8, shufB *[16]uint8, lenA uint8)