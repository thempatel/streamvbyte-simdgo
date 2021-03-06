// +build amd64

package writer

import (
	"github.com/theMPatel/streamvbyte-simdgo/pkg/encode"
	"github.com/theMPatel/streamvbyte-simdgo/pkg/shared"
)

// WriteAllFast will encode all the integers from in using the Stream VByte
// format using special hardware instructions and will return the byte array
// holding the encoded data.
func WriteAllFast(in []uint32) []byte {
	var (
		count   = len(in)
		ctrlLen = (count + 3) / 4
		stream  = make([]byte, ctrlLen+(encode.MaxBytesPerNum*count))

		dataPos  = ctrlLen
		ctrlPos  = 0
		encoded  = 0
		lowest32 = ((ctrlLen - 3) * 4) &^ 31
	)

	for ; encoded < lowest32; encoded += 32 {
		ctrls := stream[ctrlPos : ctrlPos+8]
		nums := in[encoded : encoded+32]
		out := stream[dataPos:]

		ctrl := encode.Put8uint32FastAsm(
			nums[0:8],
			out,
			shared.EncodeShuffleTable,
			shared.PerControlLenTable,
		)

		ctrls[0] = uint8(ctrl & 0xff)
		ctrls[1] = uint8(ctrl >> 8)
		sizeA := shared.ControlByteToSizeTwo(ctrl)

		ctrl = encode.Put8uint32FastAsm(
			nums[8:16],
			out[sizeA:],
			shared.EncodeShuffleTable,
			shared.PerControlLenTable,
		)

		ctrls[2] = uint8(ctrl & 0xff)
		ctrls[3] = uint8(ctrl >> 8)
		sizeB := shared.ControlByteToSizeTwo(ctrl)

		ctrl = encode.Put8uint32FastAsm(
			nums[16:24],
			out[sizeA+sizeB:],
			shared.EncodeShuffleTable,
			shared.PerControlLenTable,
		)

		ctrls[4] = uint8(ctrl & 0xff)
		ctrls[5] = uint8(ctrl >> 8)
		sizeC := shared.ControlByteToSizeTwo(ctrl)

		ctrl = encode.Put8uint32FastAsm(
			nums[24:],
			out[sizeA+sizeB+sizeC:],
			shared.EncodeShuffleTable,
			shared.PerControlLenTable,
		)

		ctrls[6] = uint8(ctrl & 0xff)
		ctrls[7] = uint8(ctrl >> 8)
		sizeD := shared.ControlByteToSizeTwo(ctrl)

		ctrlPos += 8
		dataPos += sizeA + sizeB + sizeC + sizeD
	}

	for ; ctrlPos < ctrlLen-2; ctrlPos += 2 {
		ctrl := encode.Put8uint32FastAsm(
			in[encoded:],
			stream[dataPos:],
			shared.EncodeShuffleTable,
			shared.PerControlLenTable,
		)

		stream[ctrlPos] = uint8(ctrl & 0xff)
		stream[ctrlPos+1] = uint8(ctrl >> 8)
		encoded += 8
		dataPos += shared.ControlByteToSizeTwo(ctrl)
	}

	for ; ctrlPos < ctrlLen; ctrlPos += 1 {
		nums := count - encoded
		if nums > 4 {
			nums = 4
		}
		ctrl := encode.PutUint32Scalar(in[encoded:], stream[dataPos:], nums)
		size := shared.ControlByteToSize(ctrl)
		stream[ctrlPos] = ctrl
		size -= 4 - nums
		dataPos += size
		encoded += nums
	}

	return stream[:dataPos]
}

// WriteAllDeltaFast will differentially encode all the integers from in using
// the Stream VByte format using special hardware instructions and will return
// the byte array holding the encoded data.
func WriteAllDeltaFast(in []uint32, prev uint32) []byte {
	var (
		count   = len(in)
		ctrlLen = (count + 3) / 4
		stream  = make([]byte, ctrlLen+(encode.MaxBytesPerNum*count))

		dataPos  = ctrlLen
		ctrlPos  = 0
		encoded  = 0
		lowest32 = ((ctrlLen - 3) * 4) &^ 31
	)

	for ; encoded < lowest32; encoded += 32 {
		ctrls := stream[ctrlPos : ctrlPos+8]
		nums := in[encoded : encoded+32]
		out := stream[dataPos:]

		ctrl := encode.Put8uint32DeltaFastAsm(
			nums[0:8],
			out,
			prev,
			shared.EncodeShuffleTable,
			shared.PerControlLenTable,
		)

		ctrls[0] = uint8(ctrl & 0xff)
		ctrls[1] = uint8(ctrl >> 8)
		sizeA := shared.ControlByteToSizeTwo(ctrl)

		ctrl = encode.Put8uint32DeltaFastAsm(
			nums[8:16],
			out[sizeA:],
			nums[7],
			shared.EncodeShuffleTable,
			shared.PerControlLenTable,
		)

		ctrls[2] = uint8(ctrl & 0xff)
		ctrls[3] = uint8(ctrl >> 8)
		sizeB := shared.ControlByteToSizeTwo(ctrl)

		ctrl = encode.Put8uint32DeltaFastAsm(
			nums[16:24],
			out[sizeA+sizeB:],
			nums[15],
			shared.EncodeShuffleTable,
			shared.PerControlLenTable,
		)

		ctrls[4] = uint8(ctrl & 0xff)
		ctrls[5] = uint8(ctrl >> 8)
		sizeC := shared.ControlByteToSizeTwo(ctrl)

		ctrl = encode.Put8uint32DeltaFastAsm(
			nums[24:],
			out[sizeA+sizeB+sizeC:],
			nums[23],
			shared.EncodeShuffleTable,
			shared.PerControlLenTable,
		)

		ctrls[6] = uint8(ctrl & 0xff)
		ctrls[7] = uint8(ctrl >> 8)
		sizeD := shared.ControlByteToSizeTwo(ctrl)

		ctrlPos += 8
		dataPos += sizeA + sizeB + sizeC + sizeD
		prev = nums[31]
	}

	for ; ctrlPos < ctrlLen-2; ctrlPos += 2 {
		ctrl := encode.Put8uint32DeltaFastAsm(
			in[encoded:],
			stream[dataPos:],
			prev,
			shared.EncodeShuffleTable,
			shared.PerControlLenTable,
		)

		stream[ctrlPos] = uint8(ctrl & 0xff)
		stream[ctrlPos+1] = uint8(ctrl >> 8)
		encoded += 8
		dataPos += shared.ControlByteToSizeTwo(ctrl)
		prev = in[encoded-1]
	}

	for ; ctrlPos < ctrlLen; ctrlPos += 1 {
		nums := count - encoded
		if nums > 4 {
			nums = 4
		}
		ctrl := encode.PutUint32DeltaScalar(in[encoded:], stream[dataPos:], nums, prev)
		size := shared.ControlByteToSize(ctrl)
		stream[ctrlPos] = ctrl
		size -= 4 - nums
		dataPos += size
		encoded += nums
		prev = in[encoded-1]
	}

	return stream[:dataPos]
}
