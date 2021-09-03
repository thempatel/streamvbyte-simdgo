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
	putDiffImpl Put8DiffImpl
)

type Put8Impl func(in []uint32, out []byte) (ctrl uint16)
type Put8DiffImpl func(in []uint32, out []byte, prev uint32) (ctrl uint16)

func init() {
	if GetMode() == shared.Fast {
		putImpl = Put8uint32Fast
		putDiffImpl = Put8uint32DiffFast
	} else {
		putImpl = Put8uint32Scalar
		putDiffImpl = Put8uint32DiffScalar
	}
}

// Put8uint32 is a general func you can use to encode 8 uint32's at a time.
// It will use the fastest implementation available determined during
// package initialization. If your CPU supports special hardware instructions
// then it will use an accelerated version of Stream VByte. Otherwise, the
// scalar implementation will be used as the fallback.
func Put8uint32(in []uint32, out []byte) uint16 {
	return putImpl(in, out)
}

// Put8uint32Diff is a general func you can use to encode 8 differentially coded
// uint32's with at a time. It will use the fastest implementation available
// determined during package initialization. If your CPU supports special hardware
// instructions then it will use an accelerated version of Stream VByte. Otherwise,
// the scalar implementation will be used as the fallback.
func Put8uint32Diff(in []uint32, out []byte, prev uint32) uint16 {
	return putDiffImpl(in, out, prev)
}

// Put8uint32Scalar will encode 8 uint32 values from in into out using the
// Stream VByte format. Returns an 16-bit control value produced from the
// encoding.
//
// Note: It is your responsibility to ensure that the incoming slices have
// the appropriate sizes and data otherwise this func will panic.
func Put8uint32Scalar(in []uint32, out []byte) uint16 {
	var ctrl uint16
	first := Put4uint32Scalar(in, out)
	ctrl |= uint16(first)
	encoded := shared.ControlByteToSize(first)
	second := Put4uint32Scalar(in[4:], out[encoded:])
	return ctrl | uint16(second)<<8
}

// Put4uint32Scalar will encode 4 uint32 values from in into out using the
// Stream VByte format. Returns an 8-bit control value produced from the
// encoding. Every incoming number is variably encoded, and an 8-bit control
// is constructed from the 2-bit len of each uint32. Below is an example of
// 4 uint32's and how they are encoded.
//
// 00000000 00000000 00000000 01101111  =        111
// 00000000 00000000 00000100 11010010  =       1234
// 00000000 00001100 00001010 10000011  =     789123
// 01000000 00000000 00000000 00000000  = 1073741824
//
// Num         Len      2-bit control
// ----------------------------------
// 111          1                0b00
// 1234         2                0b01
// 789123       3                0b10
// 1073741824   4                0b11
//
// Final Control byte
// 0b11100100
//
// Encoded data (little endian right-to-left bottom-to-top)
// 0b01000000 0b00000000 0b00000000 0b00000000 0b00001100
// 0b00001010 0b10000011 0b00000100 0b11010010 0b01101111
//
// Note: It is your responsibility to ensure that the incoming slices have
// the appropriate sizes and data otherwise this func will panic.
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

// Put8uint32DiffScalar will differentially encode 8 uint32 values from in into out.
// Prev provides a way for you to indicate the base value for this batch of 8.
// For example, when encoding the second batch of 8 integers out of, e.g. 16, you would
// provide a prev value of the last value in the first batch of 8 you encoded. This
// is done to ensure that the integers are correctly resolved to the correct diff. An
// example below. Note that this func assumes that the input integers are already sorted.
//
// Input:	[ 10, 20, 30, 40, 50, 60, 70, 80 ] [ 90, 100, 110, 120, 130, 140, 150, 160 ]
// Output:	[ 10, 10, 10, 10, 10, 10, 10, 10 ] [ 10, 10, 10, 10, 10, 10, 10, 10 ]
// Prev: 80
func Put8uint32DiffScalar(in []uint32, out []byte, prev uint32) uint16 {
	var ctrl uint16
	first := Put4uint32DiffScalar(in, out, prev)
	ctrl |= uint16(first)
	encoded := shared.ControlByteToSize(first)
	second := Put4uint32DiffScalar(in[4:], out[encoded:], in[3])
	return ctrl | uint16(second)<<8
}

// Put4uint32DiffScalar will differentially encode 4 uint32 values from in into out.
// Prev provides a way for you to indicate the base value for this batch of 4.
// For example, when encoding the second batch of 4 integers out of, e.g. 8, you would
// provide a prev value of the last value in the first batch of 4 you encoded. This
// is done to ensure that the integers are correctly resolved to the correct diff. An
// example below. Note that this func assumes that the input integers are already sorted.
//
// Input:	[ 10, 20, 30, 40 ] [ 50, 60, 70, 80 ]
// Output:	[ 10, 10, 10, 10 ] [ 10, 10, 10, 10 ]
// Prev: 40
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
