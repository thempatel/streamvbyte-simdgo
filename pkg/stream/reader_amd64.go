// +build amd64

package stream

import (
	_ "unsafe"

	"github.com/theMPatel/streamvbyte-simdgo/pkg/decode"
	"github.com/theMPatel/streamvbyte-simdgo/pkg/shared"
)

func AllocSlice(count int) []uint32 {
	return make([]uint32, (count + 7) & (-8))
}

func FastReadAll(count int, stream []byte, out []uint32) []uint32 {
	var (
		ctrlPos = 0
		decoded = 0
		dataPos = (count+3)/4
		ctrlLen = dataPos
		// lowest32 is the limit for the count of integers we'll read in
		// bulk 8 at a time directly from the input stream. We subtract 3
		// here since we load 16 bytes at a time in the assembly code. If
		// you attempt to load the last few control bytes worth of data,
		// it's possible there won't be enough bytes in the data stream to
		// support it, which can lead to loading from uninitialized memory.
		//
		// [ _ _ _ _ | _ _ _ _ | _ _ _ _ | _ _ _ _ ]
		//
		// Imagine the last group in the above array is all encoded with 4 bytes.
		// Decoding the first 4 integers in that group will work fine, since it
		// will load the last 3 (unused) bytes. However, when attempting to
		// decode the last three groups of 4, each load will need an extra
		// 1, 2, or 3 bytes (respectively) in order to be considered safe.
		lowest32 = ((ctrlLen-3)*4) &^ 31
	)

	for ; decoded < lowest32; decoded += 32 {
		data := stream[dataPos:]

		ctrl := uint16(stream[ctrlPos]) | uint16(stream[ctrlPos+1]) << 8
		get8uint32Fast(data, out[decoded:], ctrl)
		sizeA := shared.ControlByteToSize(stream[ctrlPos]) + shared.ControlByteToSize(stream[ctrlPos+1])

		ctrl = uint16(stream[ctrlPos+2]) | uint16(stream[ctrlPos+3]) << 8
		get8uint32Fast(data[sizeA:], out[decoded+8:], ctrl)
		sizeB := shared.ControlByteToSize(stream[ctrlPos+2]) + shared.ControlByteToSize(stream[ctrlPos+3])

		ctrl = uint16(stream[ctrlPos+4]) | uint16(stream[ctrlPos+5]) << 8
		get8uint32Fast(data[sizeA+sizeB:], out[decoded+16:], ctrl)
		sizeC := shared.ControlByteToSize(stream[ctrlPos+4]) + shared.ControlByteToSize(stream[ctrlPos+5])

		ctrl = uint16(stream[ctrlPos+6]) | uint16(stream[ctrlPos+7]) << 8
		get8uint32Fast(data[sizeA+sizeB+sizeC:], out[decoded+24:], ctrl)
		sizeD := shared.ControlByteToSize(stream[ctrlPos+6]) + shared.ControlByteToSize(stream[ctrlPos+7])

		dataPos += sizeA + sizeB + sizeC + sizeD
		ctrlPos += 8
	}

	for ; ctrlPos < ctrlLen-3; ctrlPos += 2 {
		ctrl := uint16(stream[ctrlPos]) | uint16(stream[ctrlPos+1]) << 8
		get8uint32Fast(stream[dataPos:], out[decoded:], ctrl)
		dataPos += shared.ControlByteToSize(stream[ctrlPos]) + shared.ControlByteToSize(stream[ctrlPos+1])
		decoded += 8
	}

	for ; ctrlPos < ctrlLen; ctrlPos += 1 {
		nums := count-decoded
		if nums > 4 {
			nums = 4
		}
		dataPos += decode.GetUint32Scalar(
			stream[dataPos:], out[decoded:],
			stream[ctrlPos], nums,
		)
		decoded += nums
	}

	return out[:count]
}

// get8uint32Fast here is a convenience wrapper around the assembly impl
// defined at decode.Get8uint32FastAsm. This is done to ensure that the
// assembly func is inlined at the appropriate locations. Using
// decode.Get8uint32Fast here adds another non-lined func despite the simplicity
// of the wrapper in that package.
func get8uint32Fast(in []byte, out []uint32, ctrl uint16) {
	decode.Get8uint32FastAsm(in, out, ctrl,
		shared.DecodeShuffleTable,
		shared.PerControlLenTable,
	)
}
