package writer

import (
	"github.com/theMPatel/streamvbyte-simdgo/pkg/encode"
	"github.com/theMPatel/streamvbyte-simdgo/pkg/shared"
)

func WriteAllScalar(in []uint32) []byte {
	var (
		count = len(in)
		ctrlLen = (count+3)/4
		stream = make([]byte, ctrlLen+(encode.MaxBytesPerNum*count))

		dataPos = ctrlLen
		ctrlPos = 0
		numsPos = 0
		lowest4 = count &^ 3
	)

	for ; numsPos < lowest4; numsPos += 4 {
		ctrl := encode.Put4uint32Scalar(in[numsPos:], stream[dataPos:])
		stream[ctrlPos] = ctrl
		size := shared.ControlByteToSize(ctrl)
		dataPos += size
		ctrlPos++
	}

	if lowest4 != count {
		nums := count-lowest4
		ctrl := encode.PutUint32Scalar(in[numsPos:], stream[dataPos:], nums)
		size := shared.ControlByteToSize(ctrl)
		size -= 4 - nums
		dataPos += size
		stream[ctrlPos] = ctrl
	}

	return stream[:dataPos]
}
