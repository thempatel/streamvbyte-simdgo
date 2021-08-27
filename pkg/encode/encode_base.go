// +build !amd64

package encode

import (
	"github.com/theMPatel/streamvbyte-simdgo/pkg/shared"
)

func GetMode() shared.PerformanceMode {
	return shared.Normal
}

func put8uint32(in []uint32, out []byte) uint16 {
	panic("unreachable")
}
