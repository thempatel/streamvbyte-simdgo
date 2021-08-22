package decode

import (
	"encoding/binary"
	"sync"

	"github.com/theMPatel/streamvbyte-simdgo/pkg/shared"
)

var (
	getImpl Get4Impl
	bufPool = sync.Pool{
		New: func() interface{} {
			buf := make([]byte, 4)
			return &buf
		},
	}
)

type Get4Impl func(in []byte, out []uint32, ctrl uint8) int
type cpuCheck func() bool

func init() {
	if check() {
		getImpl = get8uint32
	} else {
		getImpl = Get4uint32Scalar
	}
}

func Get4uint32(in []byte, out []uint32, ctrl uint8) int {
	return getImpl(in, out, ctrl)
}

func Get4uint32Scalar(in []byte, out []uint32, ctrl uint8) int {
	sizes := shared.PerNumLenTable[ctrl]
	n := 0
	buf := *(bufPool.Get().(*[]byte))
	defer bufPool.Put(&buf)

	len0 := sizes[0]
	len1 := sizes[1]
	len2 := sizes[2]
	len3 := sizes[3]

	out[n] = decodeOne(in, len0, buf)
	n++
	out[n] = decodeOne(in[len0:], len1, buf)
	n++
	out[n] = decodeOne(in[len0+len1:], len2, buf)
	n++
	out[n] = decodeOne(in[len0+len1+len2:], len3, buf)
	n++

	return int(len0+len1+len2+len3)
}

var zeroSlice = make([]byte, 4)

func decodeOne(input []byte, size uint8, buf []byte) uint32 {
	copy(buf, input[:size])
	copy(buf[size:], zeroSlice)
	return binary.LittleEndian.Uint32(buf)
}
