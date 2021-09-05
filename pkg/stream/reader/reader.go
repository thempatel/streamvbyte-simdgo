package reader

import (
	"github.com/theMPatel/streamvbyte-simdgo/pkg/decode"
	"github.com/theMPatel/streamvbyte-simdgo/pkg/shared"
)

func ReadAllScalar(count int, stream []byte, out []uint32) {
	var (
		ctrlLen = (count + 3) / 4

		dataPos = ctrlLen
		ctrlPos = 0
		numsPos = 0
		lowest4 = count &^ 3
	)

	for ; numsPos < lowest4; numsPos += 4 {
		ctrl := stream[ctrlPos]
		decode.Get4uint32Scalar(stream[dataPos:], out[numsPos:], ctrl)
		size := shared.ControlByteToSize(ctrl)
		dataPos += size
		ctrlPos++
	}

	if lowest4 != count {
		decode.GetUint32Scalar(stream[dataPos:], out[numsPos:], stream[ctrlPos], count-lowest4)
	}
}
