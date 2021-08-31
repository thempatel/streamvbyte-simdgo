package encode

import (
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

func Put8uint32Scalar(in []uint32, out []byte) uint16 {
	var ctrl uint16
	first := Put4uint32Scalar(in, out)
	ctrl |= uint16(first)
	encoded := shared.ControlByteToSize(first)
	second := Put4uint32Scalar(in[4:], out[encoded:])
	return ctrl | uint16(second)<<8
}

func Put4uint32Scalar(in []uint32, out []byte) uint8 {
	// bounds check hint to compiler
	_ = in[3]

	num0 := in[0]
	num1 := in[1]
	num2 := in[2]
	num3 := in[3]

	len0 := encodeOne(num0, out)
	len1 := encodeOne(num1, out[len0:])
	len2 := encodeOne(num2, out[len0+len1:])
	len3 := encodeOne(num3, out[len0+len1+len2:])

	return uint8((len0 - 1) | (len1-1)<<2 | (len2-1)<<4 | (len3-1)<<6)
}

func Put8uint32DiffScalar(in []uint32, out []byte, prev uint32) uint16 {
	var ctrl uint16
	first := Put4uint32DiffScalar(in, out, prev)
	ctrl |= uint16(first)
	encoded := shared.ControlByteToSize(first)
	second := Put4uint32DiffScalar(in[4:], out[encoded:], in[3])
	return ctrl | uint16(second)<<8
}

func Put4uint32DiffScalar(in []uint32, out []byte, prev uint32) uint8 {
	// bounds check hint to compiler
	_ = in[3]

	num0 := in[0] - prev
	num1 := in[1] - in[0]
	num2 := in[2] - in[1]
	num3 := in[3] - in[2]

	len0 := encodeOne(num0, out)
	len1 := encodeOne(num1, out[len0:])
	len2 := encodeOne(num2, out[len0+len1:])
	len3 := encodeOne(num3, out[len0+len1+len2:])

	return uint8((len0 - 1) | (len1-1)<<2 | (len2-1)<<4 | (len3-1)<<6)
}

func encodeOne(num uint32, out []byte) int {
	size := max(1, 4-(bits.LeadingZeros32(num)/8))
	switch size {
	case 4:
		out[3] = byte(num >> 24)
		fallthrough
	case 3:
		out[2] = byte(num >> 16)
		fallthrough
	case 2:
		out[1] = byte(num >> 8)
		fallthrough
	case 1:
		out[0] = byte(num)
	}
	return size
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}
