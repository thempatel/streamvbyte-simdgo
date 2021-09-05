package writer

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/theMPatel/streamvbyte-simdgo/pkg/encode"
	"github.com/theMPatel/streamvbyte-simdgo/pkg/stream/reader"
	"github.com/theMPatel/streamvbyte-simdgo/pkg/util"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestWriteAllScalar(t *testing.T) {
	for i := 0; i < 6; i++ {
		count := int(util.RandUint32() % 1e6)
		nums := util.GenUint32(count)
		stream := WriteAllScalar(nums)
		t.Run(fmt.Sprintf("WriteAll: %d", count), func(t *testing.T) {
			out := make([]uint32, count)
			reader.ReadAllScalar(count, stream, out)
			if !reflect.DeepEqual(nums, out) {
				t.Fatalf("decoded wrong nums")
			}
		})
	}
}

func TestWriteAllFast(t *testing.T) {
	for i := 0; i < 6; i++ {
		count := int(util.RandUint32() % 1e6)
		nums := util.GenUint32(count)
		stream := WriteAllScalar(nums)
		t.Run(fmt.Sprintf("WriteAll: %d", count), func(t *testing.T) {
			actual := WriteAllFast(nums)
			if !reflect.DeepEqual(stream, actual) {
				t.Fatalf("bad encoding")
			}
		})
	}
}

var readSinkA []byte

func BenchmarkWriteAllFast(b *testing.B) {
	for i := 0; i < 8; i++ {
		count := int(math.Pow10(i))
		nums := util.GenUint32(count)
		b.Run(fmt.Sprintf("Count_1e%d", i), func(b *testing.B) {
			var stream []byte
			b.SetBytes(int64(count * encode.MaxBytesPerNum))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				stream = WriteAllFast(nums)
			}
			readSinkA = stream
		})
	}
}

var readSinkB []byte

func BenchmarkFastWrite(b *testing.B) {
	count := 4096
	nums := util.GenUint32(count)
	per := count * encode.MaxBytesPerNum
	var stream []byte

	b.SetBytes(int64(per))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stream = WriteAllFast(nums)
	}
	readSinkB = stream
}

var readSinkC []byte

func BenchmarkWriteAllScalar(b *testing.B) {
	for i := 0; i < 8; i++ {
		count := int(math.Pow10(i))
		nums := util.GenUint32(count)
		b.Run(fmt.Sprintf("Count_1e%d", i), func(b *testing.B) {
			var stream []byte
			b.SetBytes(int64(count * encode.MaxBytesPerNum))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				stream = WriteAllScalar(nums)
			}
			readSinkC = stream
		})
	}
}
