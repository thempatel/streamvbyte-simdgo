// +build amd64

package reader

import (
	"github.com/theMPatel/streamvbyte-simdgo/pkg/decode"
	"github.com/theMPatel/streamvbyte-simdgo/pkg/shared"
)

func ReadAllFast(count int, stream []byte, out []uint32) {
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
		ctrls := stream[ctrlPos:ctrlPos+8] // bounds check hint
		nums := out[decoded:decoded+32]

		ctrl := uint16(ctrls[0]) | uint16(ctrls[1]) << 8
		decode.Get8uint32FastAsm(
			data,
			nums,
			ctrl,
			shared.DecodeShuffleTable,
			shared.PerControlLenTable,
		)
		sizeA := shared.ControlByteToSize(ctrls[0]) + shared.ControlByteToSize(ctrls[1])

		ctrl = uint16(ctrls[2]) | uint16(ctrls[3]) << 8
		decode.Get8uint32FastAsm(
			data[sizeA:],
			nums[8:],
			ctrl,
			shared.DecodeShuffleTable,
			shared.PerControlLenTable,
		)
		sizeB := shared.ControlByteToSize(ctrls[2]) + shared.ControlByteToSize(ctrls[3])

		ctrl = uint16(ctrls[4]) | uint16(ctrls[5]) << 8
		decode.Get8uint32FastAsm(
			data[sizeA+sizeB:],
			nums[16:],
			ctrl,
			shared.DecodeShuffleTable,
			shared.PerControlLenTable,
		)
		sizeC := shared.ControlByteToSize(ctrls[4]) + shared.ControlByteToSize(ctrls[5])

		ctrl = uint16(ctrls[6]) | uint16(ctrls[7]) << 8
		decode.Get8uint32FastAsm(
			data[sizeA+sizeB+sizeC:],
			nums[24:],
			ctrl,
			shared.DecodeShuffleTable,
			shared.PerControlLenTable,
		)
		sizeD := shared.ControlByteToSize(ctrls[6]) + shared.ControlByteToSize(ctrls[7])

		dataPos += sizeA + sizeB + sizeC + sizeD
		ctrlPos += 8
	}

	// Must be strictly less than the last 4 blocks of integers, since we can't safely
	// decode 8 if our ctrl pos starts at the first 4 in the block.
	for ; ctrlPos < ctrlLen-4; ctrlPos += 2 {
		ctrl := uint16(stream[ctrlPos]) | uint16(stream[ctrlPos+1]) << 8
		decode.Get8uint32FastAsm(
			stream[dataPos:],
			out[decoded:],
			ctrl,
			shared.DecodeShuffleTable,
			shared.PerControlLenTable,
		)
		dataPos += shared.ControlByteToSize(stream[ctrlPos]) + shared.ControlByteToSize(stream[ctrlPos+1])
		decoded += 8
	}

	for ; ctrlPos < ctrlLen; ctrlPos += 1 {
		nums := count-decoded
		if nums > 4 {
			nums = 4
		}
		dataPos += decode.GetUint32Scalar(
			stream[dataPos:],
			out[decoded:],
			stream[ctrlPos],
			nums,
		)
		decoded += nums
	}
}
