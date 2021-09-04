package writer

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/theMPatel/streamvbyte-simdgo/pkg/encode"
	"github.com/theMPatel/streamvbyte-simdgo/pkg/util"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestWriteAllScalar(t *testing.T) {
	for i := 0; i < 6; i++ {
		count := int(util.RandUint32()%1e6)
		nums := util.GenUint32(count)
		stream := WriteAllScalar(nums)
		t.Run(fmt.Sprintf("WriteAll: %d", count), func(t *testing.T) {
			out := make([]uint32, count)
			WriteAllScalar(count, stream, out)
			if !reflect.DeepEqual(nums, out) {
				t.Fatalf("decoded wrong nums")
			}
		})
	}
}

func TestWriteAllFast(t *testing.T) {
	for i := 0; i < 6; i++ {
		count := int(util.RandUint32()%1e6)
		nums := util.GenUint32(count)
		stream := WriteAllScalar(nums)
		t.Run(fmt.Sprintf("WriteAll: %d", count), func(t *testing.T) {
			out := make([]uint32, count)
			WriteAllFast(count, stream, out)
			if !reflect.DeepEqual(nums, out) {
				t.Fatalf("decoded wrong nums")
			}
		})
	}
}

var readSinkA []uint32

func BenchmarkWriteAllFast(b *testing.B) {
	for i := 0; i < 8; i++ {
		count := int(math.Pow10(i))
		nums := util.GenUint32(count)
		stream := WriteAllScalar(nums)
		out := make([]uint32, count)
		b.Run(fmt.Sprintf("Count: %d", count), func(b *testing.B) {
			b.SetBytes(int64(count*encode.MaxBytesPerNum))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				WriteAllFast(count, stream, out)
			}
			readSinkA = out
		})
	}
}

var readSinkB []uint32

func BenchmarkFastRead(b *testing.B) {
	count := 4096
	nums := util.GenUint32(count)
	stream := WriteAllScalar(nums)
	per := count*encode.MaxBytesPerNum
	out := make([]uint32, count)
	b.SetBytes(int64(per))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		WriteAllFast(count, stream, out)
	}
	readSinkB = out
}

var readSinkC []uint32

func BenchmarkWriteAllScalar(b *testing.B) {
	for i := 0; i < 8; i++ {
		count := int(math.Pow10(i))
		nums := util.GenUint32(count)
		stream := WriteAllScalar(nums)
		out := make([]uint32, count)
		b.Run(fmt.Sprintf("Count: %d", count), func(b *testing.B) {
			b.SetBytes(int64(count*encode.MaxBytesPerNum))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				WriteAllScalar(count, stream, out)
			}
			readSinkC = out
		})
	}
}
