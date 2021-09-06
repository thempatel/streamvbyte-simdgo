package writer

import (
	"github.com/theMPatel/streamvbyte-simdgo/pkg/encode"
	"github.com/theMPatel/streamvbyte-simdgo/pkg/shared"
)

const (
	jump     = 16
	jumpCtrl = jump / 4
)

// WriteAll will encode all the integers from in using the Stream VByte
// format and will return the byte array holding the encoded data. It will
// select the best implementation depending on the presence of special
// hardware instructions.
func WriteAll(in []uint32) []byte {
	if encode.GetMode() == shared.Fast {
		return WriteAllFast(in)
	} else {
		return WriteAllScalar(in)
	}
}

// WriteAllDelta will differentially encode all the integers from in using
// the Stream VByte format and will return the byte array holding the encoded
// data. It will select the best implementation depending on the presence of
// special hardware instructions.
func WriteAllDelta(in []uint32, prev uint32) []byte {
	if encode.GetMode() == shared.Fast {
		return WriteAllDeltaFast(in, prev)
	} else {
		return WriteAllDeltaScalar(in, prev)
	}
}

// WriteAllScalar will encode all the integers from in using the Stream VByte
// format and will return the byte array holding the encoded data.
func WriteAllScalar(in []uint32) []byte {
	var (
		count   = len(in)
		ctrlLen = (count + 3) / 4
		stream  = make([]byte, ctrlLen+(encode.MaxBytesPerNum*count))

		dataPos    = ctrlLen
		ctrlPos    = 0
		encoded    = 0
		lowestJump = count &^ (jump - 1)
		lowest4    = count &^ 3
	)

	for ; encoded < lowestJump; encoded += jump {
		nums := in[encoded : encoded+jump]
		data := stream[dataPos:]
		ctrls := stream[ctrlPos : ctrlPos+jumpCtrl]

		ctrl := encode.Put4uint32Scalar(nums, data)
		ctrls[0] = ctrl
		sizeA := shared.ControlByteToSize(ctrl)

		ctrl = encode.Put4uint32Scalar(nums[4:], data[sizeA:])
		ctrls[1] = ctrl
		sizeB := shared.ControlByteToSize(ctrl)

		ctrl = encode.Put4uint32Scalar(nums[8:], data[sizeA+sizeB:])
		ctrls[2] = ctrl
		sizeC := shared.ControlByteToSize(ctrl)

		ctrl = encode.Put4uint32Scalar(nums[12:], data[sizeA+sizeB+sizeC:])
		ctrls[3] = ctrl
		sizeD := shared.ControlByteToSize(ctrl)

		dataPos += sizeA + sizeB + sizeC + sizeD
		ctrlPos += jumpCtrl
	}

	for ; encoded < lowest4; encoded += 4 {
		ctrl := encode.Put4uint32Scalar(in[encoded:], stream[dataPos:])
		stream[ctrlPos] = ctrl
		size := shared.ControlByteToSize(ctrl)
		dataPos += size
		ctrlPos++
	}

	if lowest4 != count {
		nums := count - lowest4
		ctrl := encode.PutUint32Scalar(in[encoded:], stream[dataPos:], nums)
		size := shared.ControlByteToSize(ctrl)
		size -= 4 - nums
		dataPos += size
		stream[ctrlPos] = ctrl
	}

	return stream[:dataPos]
}

// WriteAllDeltaScalar will differentially encode all the integers from in using
// the Stream VByte format and will return the byte array holding the encoded data.
func WriteAllDeltaScalar(in []uint32, prev uint32) []byte {
	var (
		count   = len(in)
		ctrlLen = (count + 3) / 4
		stream  = make([]byte, ctrlLen+(encode.MaxBytesPerNum*count))

		dataPos    = ctrlLen
		ctrlPos    = 0
		encoded    = 0
		lowestJump = count &^ (jump - 1)
		lowest4    = count &^ 3
	)

	for ; encoded < lowestJump; encoded += jump {
		nums := in[encoded : encoded+jump]
		data := stream[dataPos:]
		ctrls := stream[ctrlPos : ctrlPos+jumpCtrl]

		ctrl := encode.Put4uint32DeltaScalar(nums, data, prev)
		ctrls[0] = ctrl
		sizeA := shared.ControlByteToSize(ctrl)

		ctrl = encode.Put4uint32DeltaScalar(nums[4:], data[sizeA:], nums[3])
		ctrls[1] = ctrl
		sizeB := shared.ControlByteToSize(ctrl)

		ctrl = encode.Put4uint32DeltaScalar(nums[8:], data[sizeA+sizeB:], nums[7])
		ctrls[2] = ctrl
		sizeC := shared.ControlByteToSize(ctrl)

		ctrl = encode.Put4uint32DeltaScalar(nums[12:], data[sizeA+sizeB+sizeC:], nums[11])
		ctrls[3] = ctrl
		sizeD := shared.ControlByteToSize(ctrl)

		dataPos += sizeA + sizeB + sizeC + sizeD
		ctrlPos += jumpCtrl
		prev = nums[15]
	}

	for ; encoded < lowest4; encoded += 4 {
		ctrl := encode.Put4uint32DeltaScalar(in[encoded:], stream[dataPos:], prev)
		stream[ctrlPos] = ctrl
		size := shared.ControlByteToSize(ctrl)
		dataPos += size
		ctrlPos++
		prev = in[encoded+3]
	}

	if lowest4 != count {
		nums := count - lowest4
		ctrl := encode.PutUint32DeltaScalar(in[encoded:], stream[dataPos:], nums, prev)
		size := shared.ControlByteToSize(ctrl)
		size -= 4 - nums
		dataPos += size
		stream[ctrlPos] = ctrl
	}

	return stream[:dataPos]
}
