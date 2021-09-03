// +build !amd64

package decode

import (
	"github.com/theMPatel/streamvbyte-simdgo/pkg/shared"
)

func GetMode() shared.PerformanceMode {
	return shared.Normal
}

func Get8uint32Fast(in []byte, out []uint32, ctrl uint16) int {
	panic("unreachable")
}

func Get8uint32DiffFast(in []byte, out []uint32, ctrl uint16, prev uint32) int {
	panic("unreachable")
}
