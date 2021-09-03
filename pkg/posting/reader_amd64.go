// +build amd64

package posting

import (
	"github.com/theMPatel/streamvbyte-simdgo/pkg/decode"
	"github.com/theMPatel/streamvbyte-simdgo/pkg/shared"
)

// fastReader is a Stream VByte reader that uses solely hardware accelerated
// algorithms to operate on an input stream.
type fastReader struct {
	stream []byte
	count, dataPos, ctrlPos int
}

func NewFastReader(count int, stream []byte) Reader {
	return &fastReader{
		stream: stream,
		count: count,
		ctrlPos: 0,
		dataPos: (count+3)/4,
	}
}

func (f *fastReader) Read(count int) []uint32 {
	out := make([]uint32, (count + 7) & (-8)) // Round up to nearest 8

	lowest8 := count &^ 7
	for n := 0; n < lowest8; n += 8 {
		ctrl := uint16(f.stream[f.ctrlPos]) | uint16(f.stream[f.ctrlPos+1]) << 8
		size := shared.ControlByteToSizeTwo(ctrl)
		decode.Get8uint32Fast(f.stream[f.dataPos:], out[n:], ctrl)
		f.dataPos += size
		f.ctrlPos += 2
	}



	return out[:count]
}
