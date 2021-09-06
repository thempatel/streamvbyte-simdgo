package writer

import (
	"encoding/binary"
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

func TestWriteAllDeltaScalar(t *testing.T) {
	for i := 0; i < 6; i++ {
		count := int(util.RandUint32() % 1e6)
		nums := util.GenUint32(count)
		util.SortUint32(nums)
		diffed := make([]uint32, count)
		util.Delta(nums, diffed)

		stream := WriteAllScalar(diffed)
		t.Run(fmt.Sprintf("WriteAll: %d", count), func(t *testing.T) {
			actual := WriteAllDeltaScalar(nums, 0)
			if !reflect.DeepEqual(stream, actual) {
				t.Fatalf("bad encoding")
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

func TestWriteAllDeltaFast(t *testing.T) {
	for i := 0; i < 6; i++ {
		count := int(util.RandUint32() % 1e6)
		nums := util.GenUint32(count)
		util.SortUint32(nums)
		diffed := make([]uint32, count)
		util.Delta(nums, diffed)

		stream := WriteAllScalar(diffed)
		t.Run(fmt.Sprintf("WriteAll: %d", count), func(t *testing.T) {
			actual := WriteAllDeltaFast(nums, 0)
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

func BenchmarkWriteAllDeltaFast(b *testing.B) {
	for i := 0; i < 8; i++ {
		count := int(math.Pow10(i))
		nums := util.GenUint32(count)
		util.SortUint32(nums)
		b.Run(fmt.Sprintf("Count_1e%d", i), func(b *testing.B) {
			var stream []byte
			b.SetBytes(int64(count * encode.MaxBytesPerNum))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				stream = WriteAllDeltaFast(nums, 0)
			}
			readSinkB = stream
		})
	}
}

var readSinkC []byte

func BenchmarkWriteAllScalar(b *testing.B) {
	for i := 0; i < 8; i++ {
		count := int(math.Pow10(i))
		nums := util.GenUint32(count)
		util.SortUint32(nums)
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

var readSinkD []byte

func BenchmarkWriteAllDeltaScalar(b *testing.B) {
	for i := 0; i < 8; i++ {
		count := int(math.Pow10(i))
		nums := util.GenUint32(count)
		util.SortUint32(nums)
		b.Run(fmt.Sprintf("Count_1e%d", i), func(b *testing.B) {
			var stream []byte
			b.SetBytes(int64(count * encode.MaxBytesPerNum))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				stream = WriteAllDeltaScalar(nums, 0)
			}
			readSinkD = stream
		})
	}
}

var readSinkE int

func BenchmarkWriteAllVarint(b *testing.B) {
	for i := 0; i < 8; i++ {
		count := int(math.Pow10(i))
		nums := util.GenUint32(count)
		util.SortUint32(nums)
		out := make([]byte, count*binary.MaxVarintLen32)
		written := 0
		b.Run(fmt.Sprintf("Count_1e%d", i), func(b *testing.B) {
			b.SetBytes(int64(count * encode.MaxBytesPerNum))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				written = util.PutVarint(nums, out)
			}
			readSinkE = written
		})
	}
}

var readSinkF int

func BenchmarkWriteAllDeltaVarint(b *testing.B) {
	for i := 0; i < 8; i++ {
		count := int(math.Pow10(i))
		nums := util.GenUint32(count)
		util.SortUint32(nums)
		out := make([]byte, count*binary.MaxVarintLen32)
		written := 0
		b.Run(fmt.Sprintf("Count_1e%d", i), func(b *testing.B) {
			b.SetBytes(int64(count * encode.MaxBytesPerNum))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				written = util.PutDeltaVarint(nums, out, 0)
			}
			readSinkE = written
		})
	}
}
