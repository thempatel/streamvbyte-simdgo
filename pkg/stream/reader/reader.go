package reader

import (
	"github.com/theMPatel/streamvbyte-simdgo/pkg/decode"
	"github.com/theMPatel/streamvbyte-simdgo/pkg/shared"
)

const (
	jump     = 16
	jumpCtrl = jump / 4
)

// ReadAll will read the entire input stream into out according to the
// Stream VByte format. It will select the best implementation depending
// on the presence of special hardware instructions.
//
// Note: It is your responsibility to ensure that the incoming slices are
// appropriately sized as well as tracking the count of integers in the
// stream.
func ReadAll(count int, stream []byte, out []uint32) {
	if decode.GetMode() == shared.Fast {
		ReadAllFast(count, stream, out)
	} else {
		ReadAllScalar(count, stream, out)
	}
}

// ReadAllScalar will read the entire input stream into out according to the
// Stream VByte format.
//
// Note: It is your responsibility to ensure that the incoming slices are
// appropriately sized as well as tracking the count of integers in the
// stream.
func ReadAllScalar(count int, stream []byte, out []uint32) {
	var (
		ctrlLen = (count + 3) / 4

		dataPos    = ctrlLen
		ctrlPos    = 0
		decoded    = 0
		lowestJump = count &^ (jump - 1)
		lowest4    = count &^ 3
	)

	for ; decoded < lowestJump; decoded += jump {
		data := stream[dataPos:]
		ctrls := stream[ctrlPos : ctrlPos+jumpCtrl]
		nums := out[decoded : decoded+jump]

		ctrl := ctrls[0]
		decode.Get4uint32Scalar(data, nums, ctrl)
		sizeA := shared.ControlByteToSize(ctrl)

		ctrl = ctrls[1]
		decode.Get4uint32Scalar(data[sizeA:], nums[4:], ctrl)
		sizeB := shared.ControlByteToSize(ctrl)

		ctrl = ctrls[2]
		decode.Get4uint32Scalar(data[sizeA+sizeB:], nums[8:], ctrl)
		sizeC := shared.ControlByteToSize(ctrl)

		ctrl = ctrls[3]
		decode.Get4uint32Scalar(data[sizeA+sizeB+sizeC:], nums[12:], ctrl)
		sizeD := shared.ControlByteToSize(ctrl)

		dataPos += sizeA + sizeB + sizeC + sizeD
		ctrlPos += jumpCtrl
	}

	for ; decoded < lowest4; decoded += 4 {
		ctrl := stream[ctrlPos]
		decode.Get4uint32Scalar(stream[dataPos:], out[decoded:], ctrl)
		size := shared.ControlByteToSize(ctrl)
		dataPos += size
		ctrlPos++
	}

	if lowest4 != count {
		decode.GetUint32Scalar(stream[dataPos:], out[decoded:], stream[ctrlPos], count-lowest4)
	}
}

// ReadAllDeltaScalar will read the entire input stream into out according to the
// Stream VByte format. It will reconstruct the original non differentially
// encoded values.
//
// Note: It is your responsibility to ensure that the incoming slices are
// appropriately sized as well as tracking the count of integers in the
// stream.
func ReadAllDeltaScalar(count int, stream []byte, out []uint32, prev uint32) {
	var (
		ctrlLen = (count + 3) / 4

		dataPos    = ctrlLen
		ctrlPos    = 0
		decoded    = 0
		lowestJump = count &^ (jump - 1)
		lowest4    = count &^ 3
	)

	for ; decoded < lowestJump; decoded += jump {
		data := stream[dataPos:]
		ctrls := stream[ctrlPos : ctrlPos+jumpCtrl]
		nums := out[decoded : decoded+jump]

		ctrl := ctrls[0]
		decode.Get4uint32DeltaScalar(data, nums, ctrl, prev)
		sizeA := shared.ControlByteToSize(ctrl)

		ctrl = ctrls[1]
		decode.Get4uint32DeltaScalar(data[sizeA:], nums[4:], ctrl, nums[3])
		sizeB := shared.ControlByteToSize(ctrl)

		ctrl = ctrls[2]
		decode.Get4uint32DeltaScalar(data[sizeA+sizeB:], nums[8:], ctrl, nums[7])
		sizeC := shared.ControlByteToSize(ctrl)

		ctrl = ctrls[3]
		decode.Get4uint32DeltaScalar(data[sizeA+sizeB+sizeC:], nums[12:], ctrl, nums[11])
		sizeD := shared.ControlByteToSize(ctrl)

		dataPos += sizeA + sizeB + sizeC + sizeD
		ctrlPos += jumpCtrl
		prev = nums[15]
	}

	for ; decoded < lowest4; decoded += 4 {
		ctrl := stream[ctrlPos]
		decode.Get4uint32DeltaScalar(stream[dataPos:], out[decoded:], ctrl, prev)
		size := shared.ControlByteToSize(ctrl)
		dataPos += size
		ctrlPos++
		prev = out[decoded+3]
	}

	if lowest4 != count {
		decode.GetUint32DeltaScalar(stream[dataPos:], out[decoded:], stream[ctrlPos], count-lowest4, prev)
	}
}
