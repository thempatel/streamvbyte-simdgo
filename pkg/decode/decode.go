package decode

import (
	"github.com/theMPatel/streamvbyte-simdgo/pkg/shared"
)

var (
	getImpl Get8Impl
)

type Get8Impl func(in []byte, out []uint32, ctrl uint16) int

func init() {
	if GetMode() == shared.Fast {
		getImpl = get8uint32
	} else {
		getImpl = Get8uint32Scalar
	}
}

// Get8uint32 is a general func you can use to decode 8 uint32's at a time.
// It will use the fastest implementation available determined during
// package initialization. If your CPU supports special hardware instructions
// then it will use an accelerated version of Stream VByte. Otherwise, the
// scalar implementation will be used as the fallback.
func Get8uint32(in []byte, out []uint32, ctrl uint16) int {
	return getImpl(in, out, ctrl)
}

// Get8uint32Scalar will decode 8 uint32 values from in into out using the
// Stream VByte format. Returns the number of bytes read from the input
// buffer.
//
// Note: It is your responsibility to ensure that the incoming slices have
// the appropriate sizes and data otherwise this func will panic.
func Get8uint32Scalar(in []byte, out []uint32, ctrl uint16) int {
	read := Get4uint32Scalar(in, out, uint8(ctrl&0xff))
	return read + Get4uint32Scalar(in[read:], out[4:], uint8(ctrl>>8))
}

// Get4uint32Scalar will decode 4 uint32 values from in into out using the
// Stream VByte format. Returns the number of bytes read from the input
// buffer.
//
// Note: It is your responsibility to ensure that the incoming slices have
// the appropriate sizes and data otherwise this func will panic.
func Get4uint32Scalar(in []byte, out []uint32, ctrl uint8) int {
	sizes := shared.PerNumLenTable[ctrl]

	len0 := sizes[0]
	len1 := sizes[1]
	len2 := sizes[2]
	len3 := sizes[3]

	out[0] = decodeOne(in, len0)
	out[1] = decodeOne(in[len0:], len1)
	out[2] = decodeOne(in[len0+len1:], len2)
	out[3] = decodeOne(in[len0+len1+len2:], len3)

	return int(len0 + len1 + len2 + len3)
}

// Get8uint32DiffScalar will decode 8 uint32 values from in into out and reconstruct
// the original values via differential coding. Prev provides a way for you to
// indicate the base value for this batch of 8. For example, when decoding the second
// batch of 8 integers out of, e.g. 16, you would provide a prev value of the last value
// in the first batch of 8 you decoded. This is done to ensure that the integers are
// correctly resolved to the correct diff. An example below.
//
// Input:	[ 10, 10, 10, 10, 10, 10, 10, 10 ] [ 10, 10, 10, 10, 10, 10, 10, 10 ]
// Output:	[ 10, 20, 30, 40, 50, 60, 70, 80 ] [ 90, 100, 110, 120, 130, 140, 150, 160 ]
// Prev: 80
func Get8uint32DiffScalar(in []byte, out []uint32, ctrl uint16, prev uint32) int {
	read := Get4uint32DiffScalar(in, out, uint8(ctrl&0xff), prev)
	return read + Get4uint32DiffScalar(in[read:], out[4:], uint8(ctrl>>8), out[3])
}

// Get4uint32DiffScalar will decode 4 uint32 values from in into out and reconstruct
// the original values via differential coding. Prev provides a way for you to
// indicate the base value for this batch of 4. For example, when decoding the second
// batch of 4 integers out of, e.g. 8, you would provide a prev value of the last value
// in the first batch of 4 you decoded. This is done to ensure that the integers are
// correctly resolved to the correct diff. An example below.
//
// Input:	[ 10, 10, 10, 10 ] [ 10, 10, 10, 10 ]
// Output:	[ 10, 20, 30, 40 ] [ 50, 60, 70, 80 ]
// Prev: 40
func Get4uint32DiffScalar(in []byte, out []uint32, ctrl uint8, prev uint32) int {
	// bounds check hint to compiler
	_ = out[3]

	sizes := shared.PerNumLenTable[ctrl]

	len0 := sizes[0]
	len1 := sizes[1]
	len2 := sizes[2]
	len3 := sizes[3]

	out[0] = decodeOne(in, len0) + prev
	out[1] = decodeOne(in[len0:], len1) + out[0]
	out[2] = decodeOne(in[len0+len1:], len2) + out[1]
	out[3] = decodeOne(in[len0+len1+len2:], len3) + out[2]

	return int(len0 + len1 + len2 + len3)
}

func decodeOne(b []byte, size uint8) uint32 {
	switch size {
	case 4:
		return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
	case 3:
		return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16
	case 2:
		return uint32(b[0]) | uint32(b[1])<<8
	case 1:
		return uint32(b[0])
	}
	panic("impossible")
}