package encode

import "github.com/theMPatel/streamvbyte-simdgo/internal/encode/x86"

func ControlBytes(in []uint32) uint32 {
	return x86.ControlBytes(in)
}
