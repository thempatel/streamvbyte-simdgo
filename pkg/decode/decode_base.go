// +build !amd64

package decode

import (
	"github.com/theMPatel/streamvbyte-simdgo/pkg/shared"
)

func GetMode() shared.PerformanceMode {
	return shared.Normal
}

func get8uint32(in []byte, out []uint32, ctrl uint16) int {
	panic("unreachable")
}
