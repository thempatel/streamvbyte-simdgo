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

func Get8uint32(in []byte, out []uint32, ctrl uint16) int {
	return getImpl(in, out, ctrl)
}

func Get8uint32Scalar(in []byte, out []uint32, ctrl uint16) int {
	read := Get4uint32Scalar(in, out, uint8(ctrl&0xff))
	return read + Get4uint32Scalar(in[read:], out[4:], uint8(ctrl>>8))
}

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

func Get8uint32DiffScalar(in []byte, out []uint32, ctrl uint16, prev uint32) int {
	read := Get4uint32DiffScalar(in, out, uint8(ctrl&0xff), prev)
	return read + Get4uint32DiffScalar(in[read:], out[4:], uint8(ctrl>>8), out[3])
}

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

var dispatch = [5]func([]byte) uint32 {
	nil,
	decode1byte,
	decode2byte,
	decode3byte,
	decode4byte,
}

func decodeOne(input []byte, size uint8) uint32 {
	return dispatch[size](input)
}

func decode1byte(b []byte) uint32 {
	return uint32(b[0])
}

func decode2byte(b []byte) uint32 {
	_ = b[1] // bounds check hint to compiler
	return uint32(b[0]) | uint32(b[1])<<8
}

func decode3byte(b []byte) uint32 {
	_ = b[2] // bounds check hint to compiler
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16
}

func decode4byte(b []byte) uint32 {
	_ = b[3] // bounds check hint to compiler
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
}
