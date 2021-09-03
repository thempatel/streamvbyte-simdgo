package stream

import (
	"reflect"
	"testing"

	"github.com/theMPatel/streamvbyte-simdgo/pkg/encode"
	"github.com/theMPatel/streamvbyte-simdgo/pkg/shared"
	"github.com/theMPatel/streamvbyte-simdgo/pkg/util"
)

func encodeNums(nums []uint32) []byte {
	count := len(nums)
	data := make([]byte, encode.MaxBytesPerNum*count)
	ctrlData := make([]byte, 0, (count+3)/4)

	dataPos := 0
	numsPos := 0
	lowest4 := count &^ 3
	for ; numsPos < lowest4; numsPos += 4 {
		ctrl := encode.Put4uint32Scalar(nums[numsPos:], data[dataPos:])
		size := shared.ControlByteToSize(ctrl)
		ctrlData = append(ctrlData, ctrl)
		dataPos += size
	}

	if lowest4 != count {
		rest := make([]uint32, 4)
		copy(rest, nums[numsPos:])

		ctrl := encode.Put4uint32Scalar(rest, data[dataPos:])
		dataPos += shared.ControlByteToSize(ctrl)-(count-lowest4)
		ctrlData = append(ctrlData, ctrl)

	}

	final := make([]byte, len(ctrlData)+dataPos)
	copy(final, ctrlData)
	copy(final[len(ctrlData):], data[:dataPos])
	return final
}

func TestFastReadAll(t *testing.T) {
	count := 100
	nums := util.GenUint32(count)
	stream := encodeNums(nums)
	readNums := FastReadAll(count, stream)

	if !reflect.DeepEqual(nums, readNums) {
		t.Fatalf("decoded wrong nums")
	}
}
