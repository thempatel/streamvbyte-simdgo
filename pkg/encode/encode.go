package encode

import (
	"encoding/binary"
	"math/bits"

	"github.com/theMPatel/streamvbyte-simdgo/pkg/shared"
)

const (
	MaxBytesPerNum = 4
)

var (
	putImpl Put8Impl
)

type Put8Impl func([]uint32, []byte) uint16

func init() {
	if GetMode() == shared.Fast {
		putImpl = put8uint32
	} else {
		putImpl = Put8uint32Scalar
	}
}

func Put8uint32(in []uint32, out []byte) uint16 {
	return putImpl(in, out)
}

func Put8uint32Scalar(in []uint32, out []byte ) uint16 {
	var ctrl uint16
	first := Put4uint32Scalar(in, out)
	ctrl |= uint16(first)
	encoded := shared.ControlByteToSize(first)
	second := Put4uint32Scalar(in[4:], out[encoded:])
	return ctrl | uint16(second)<<8
}

func Put4uint32Scalar(in []uint32, out []byte ) uint8 {
	// Drop the bounds checks
	_ = in[3]

	num0 := in[0]
	num1 := in[1]
	num2 := in[2]
	num3 := in[3]

	len0 := encodeOne(num0, out)
	len1 := encodeOne(num1, out[len0:])
	len2 := encodeOne(num2, out[len0+len1:])
	len3 := encodeOne(num3, out[len0+len1+len2:])

	return uint8((len0-1) | (len1-1) << 2 | (len2-1) << 4 | (len3-1) << 6)
}

func encodeOne(num uint32, out[]byte) int {
	size := max(1, 4 - (bits.LeadingZeros32(num) / 8))
	binary.LittleEndian.PutUint32(out, num)
	return size
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}