package decode

import (
	"github.com/theMPatel/streamvbyte-simdgo/pkg/shared"
)

var (
	getImpl      Get8Impl
	getDeltaImpl Get8DeltaImpl
)

type Get8Impl func(in []byte, out []uint32, ctrl uint16)
type Get8DeltaImpl func(in []byte, out []uint32, ctrl uint16, prev uint32)

func init() {
	if GetMode() == shared.Fast {
		getImpl = Get8uint32Fast
		getDeltaImpl = Get8uint32DeltaFast
	} else {
		getImpl = Get8uint32Scalar
		getDeltaImpl = Get8uint32DeltaScalar
	}
}

// Get8uint32 is a general func you can use to decode 8 uint32's at a time.
// It will use the fastest implementation available determined during
// package initialization. If your CPU supports special hardware instructions
// then it will use an accelerated version of Stream VByte. Otherwise, the
// scalar implementation will be used as the fallback.
func Get8uint32(in []byte, out []uint32, ctrl uint16) {
	getImpl(in, out, ctrl)
}

// Get8uint32Delta is a general func you can use to decode 8 differentially coded
// uint32's at a time. It will use the fastest implementation available determined
// during package initialization. If your CPU supports special hardware instructions
// then it will use an accelerated version of Stream VByte. Otherwise, the
// scalar implementation will be used as the fallback.
func Get8uint32Delta(in []byte, out []uint32, ctrl uint16, prev uint32) {
	getDeltaImpl(in, out, ctrl, prev)
}

// Get8uint32Scalar will decode 8 uint32 values from in into out using the
// Stream VByte format. Returns the number of bytes read from the input
// buffer.
//
// Note: It is your responsibility to ensure that the incoming slices have
// the appropriate sizes and data otherwise this func will panic.
func Get8uint32Scalar(in []byte, out []uint32, ctrl uint16) {
	lower := uint8(ctrl & 0xff)
	upper := uint8(ctrl >> 8)
	lowerSize := shared.ControlByteToSize(lower)
	Get4uint32Scalar(in, out, lower)
	Get4uint32Scalar(in[lowerSize:], out[4:], upper)
}

// Get4uint32Scalar will decode 4 uint32 values from in into out using the
// Stream VByte format. Returns the number of bytes read from the input
// buffer.
//
// Note: It is your responsibility to ensure that the incoming slices have
// the appropriate sizes and data otherwise this func will panic.
func Get4uint32Scalar(in []byte, out []uint32, ctrl uint8) {
	sizes := shared.PerNumLenTable[ctrl]

	len3 := sizes[3]
	len2 := sizes[2]
	len1 := sizes[1]
	len0 := sizes[0]

	out[3] = decodeOne(in[len0+len1+len2:], len3)
	out[2] = decodeOne(in[len0+len1:], len2)
	out[1] = decodeOne(in[len0:], len1)
	out[0] = decodeOne(in, len0)
}

// GetUint32Scalar decodes up to 4 integers from in into out using the
// Stream VByte format.
//
// Note: It is your responsibility to ensure that the incoming slices have
// the appropriate sizes and data otherwise this func will panic.
func GetUint32Scalar(in []byte, out []uint32, ctrl uint8, count int) int {
	if count == 0 {
		return 0
	}

	if count > 4 {
		count = 4
	}

	shift := 0
	total := 0
	for i := 0; i < count; i++ {
		size := ((ctrl >> shift) & 0x3) + 1
		out[i] = decodeOne(in[total:], size)
		total += int(size)
		shift += 2
	}

	return total
}

// GetUint32DeltaScalar decodes up to 4 integers from in into out using the
// Stream VByte format. It will reconstruct the original non differentially
// encoded values.
//
// Note: It is your responsibility to ensure that the incoming slices have
// the appropriate sizes and data otherwise this func will panic.
func GetUint32DeltaScalar(in []byte, out []uint32, ctrl uint8, count int, prev uint32) int {
	if count == 0 {
		return 0
	}

	if count > 4 {
		count = 4
	}

	shift := 0
	total := 0
	for i := 0; i < count; i++ {
		size := ((ctrl >> shift) & 0x3) + 1
		num := decodeOne(in[total:], size) + prev
		out[i] = num
		prev = num
		total += int(size)
		shift += 2
	}

	return total
}

// Get8uint32DeltaScalar will decode 8 uint32 values from in into out and reconstruct
// the original values via differential coding. Prev provides a way for you to
// indicate the base value for this batch of 8. For example, when decoding the second
// batch of 8 integers out of, e.g. 16, you would provide a prev value of the last value
// in the first batch of 8 you decoded. This is done to ensure that the integers are
// correctly resolved to the correct diff. An example below.
//
// Input:	[ 10, 10, 10, 10, 10, 10, 10, 10 ] [ 10, 10, 10, 10, 10, 10, 10, 10 ]
// Output:	[ 10, 20, 30, 40, 50, 60, 70, 80 ] [ 90, 100, 110, 120, 130, 140, 150, 160 ]
// Prev: 80
func Get8uint32DeltaScalar(in []byte, out []uint32, ctrl uint16, prev uint32) {
	lower := uint8(ctrl & 0xff)
	upper := uint8(ctrl >> 8)
	lowerSize := shared.ControlByteToSize(lower)
	Get4uint32DeltaScalar(in, out, lower, prev)
	Get4uint32DeltaScalar(in[lowerSize:], out[4:], upper, out[3])
}

// Get4uint32DeltaScalar will decode 4 uint32 values from in into out and reconstruct
// the original values via differential coding. Prev provides a way for you to
// indicate the base value for this batch of 4. For example, when decoding the second
// batch of 4 integers out of, e.g. 8, you would provide a prev value of the last value
// in the first batch of 4 you decoded. This is done to ensure that the integers are
// correctly resolved to the correct diff. An example below.
//
// Input:	[ 10, 10, 10, 10 ] [ 10, 10, 10, 10 ]
// Output:	[ 10, 20, 30, 40 ] [ 50, 60, 70, 80 ]
// Prev: 40
func Get4uint32DeltaScalar(in []byte, out []uint32, ctrl uint8, prev uint32) {
	sizes := shared.PerNumLenTable[ctrl]

	len0 := sizes[0]
	len1 := sizes[1]
	len2 := sizes[2]
	len3 := sizes[3]

	// bounds check hint to compiler
	_ = out[3]
	out[0] = decodeOne(in, len0) + prev
	out[1] = decodeOne(in[len0:], len1) + out[0]
	out[2] = decodeOne(in[len0+len1:], len2) + out[1]
	out[3] = decodeOne(in[len0+len1+len2:], len3) + out[2]
}

func decodeOne(b []byte, size uint8) uint32 {
	switch size {
	case 4:
		return uint32(b[3])<<24 | uint32(b[2])<<16 | uint32(b[1])<<8 | uint32(b[0])
	case 3:
		return uint32(b[2])<<16 | uint32(b[1])<<8 | uint32(b[0])
	case 2:
		return uint32(b[1])<<8 | uint32(b[0])
	case 1:
		return uint32(b[0])
	}
	panic("impossible")
}
