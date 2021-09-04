package stream

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/theMPatel/streamvbyte-simdgo/pkg/encode"
	"github.com/theMPatel/streamvbyte-simdgo/pkg/shared"
	"github.com/theMPatel/streamvbyte-simdgo/pkg/util"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

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
		buf := make([]byte, encode.MaxBytesPerNum*4)
		copy(rest, nums[numsPos:])

		ctrl := encode.Put4uint32Scalar(rest, buf)
		size := shared.ControlByteToSize(ctrl)
		size -= 4 - (count - lowest4)
		copy(data[dataPos:], buf[:size])
		dataPos += size
		ctrlData = append(ctrlData, ctrl)

	}

	final := make([]byte, len(ctrlData)+dataPos)
	copy(final, ctrlData)
	copy(final[len(ctrlData):], data[:dataPos])
	return final
}

func TestReadAllFast(t *testing.T) {
	for i := 0; i < 6; i++ {
		count := int(util.RandUint32()%1e6)
		nums := util.GenUint32(count)
		stream := encodeNums(nums)
		t.Run(fmt.Sprintf("ReadAll: %d", count), func(t *testing.T) {
			out := make([]uint32, count)
			readNums := ReadAllFast(count, stream, out)
			if !reflect.DeepEqual(nums, readNums) {
				t.Fatalf("decoded wrong nums")
			}
		})
	}
}

var readSinkA []uint32

func BenchmarkReadAllFast(b *testing.B) {
	for i := 0; i < 8; i++ {
		count := int(math.Pow10(i))
		nums := util.GenUint32(count)
		stream := encodeNums(nums)
		out := make([]uint32, count)
		b.Run(fmt.Sprintf("Count: %d", count), func(b *testing.B) {
			b.SetBytes(int64(count*encode.MaxBytesPerNum))
			b.ResetTimer()
			var read []uint32
			for i := 0; i < b.N; i++ {
				read = ReadAllFast(count, stream, out)
			}
			readSinkA = read
		})
	}
}

var readSinkB []uint32

func BenchmarkFastRead(b *testing.B) {
	count := 4096
	nums := util.GenUint32(count)
	stream := encodeNums(nums)
	per := count*encode.MaxBytesPerNum
	out := make([]uint32, count)
	b.SetBytes(int64(per))
	b.ResetTimer()
	var read []uint32
	for i := 0; i < b.N; i++ {
		read = ReadAllFast(count, stream, out)
	}
	readSinkB = read
}