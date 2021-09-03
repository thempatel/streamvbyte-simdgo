// +build !amd64

package encode

import (
	"github.com/theMPatel/streamvbyte-simdgo/pkg/shared"
)

func GetMode() shared.PerformanceMode {
	return shared.Normal
}

func Put8uint32Fast(in []uint32, out []byte) uint16 {
	panic("unreachable")
}

func Put8uint32DiffFast(in []uint32, out []byte, prev uint32) uint16 {
	panic("unreachable")
}
