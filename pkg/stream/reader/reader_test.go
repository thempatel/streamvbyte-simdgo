package reader

import (
	"encoding/binary"
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/theMPatel/streamvbyte-simdgo/pkg/encode"
	"github.com/theMPatel/streamvbyte-simdgo/pkg/stream/writer"
	"github.com/theMPatel/streamvbyte-simdgo/pkg/util"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestReadAllScalar(t *testing.T) {
	for i := 0; i < 6; i++ {
		count := int(util.RandUint32() % 1e6)
		nums := util.GenUint32(count)
		stream := writer.WriteAllScalar(nums)
		t.Run(fmt.Sprintf("ReadAll: %d", count), func(t *testing.T) {
			out := make([]uint32, count)
			ReadAllScalar(count, stream, out)
			if !reflect.DeepEqual(nums, out) {
				t.Fatalf("decoded wrong nums")
			}
		})
	}
}

func TestReadAllDeltaScalar(t *testing.T) {
	for i := 0; i < 6; i++ {
		count := int(util.RandUint32() % 1e6)
		nums := util.GenUint32(count)
		util.SortUint32(nums)
		stream := writer.WriteAllDeltaScalar(nums, 0)
		t.Run(fmt.Sprintf("ReadAll: %d", count), func(t *testing.T) {
			out := make([]uint32, count)
			ReadAllDeltaScalar(count, stream, out, 0)
			if !reflect.DeepEqual(nums, out) {
				t.Fatalf("decoded wrong nums")
			}
		})
	}
}

func TestReadAllFast(t *testing.T) {
	for i := 0; i < 6; i++ {
		count := int(util.RandUint32() % 1e6)
		nums := util.GenUint32(count)
		stream := writer.WriteAllScalar(nums)
		t.Run(fmt.Sprintf("ReadAll: %d", count), func(t *testing.T) {
			out := make([]uint32, count)
			ReadAllFast(count, stream, out)
			if !reflect.DeepEqual(nums, out) {
				t.Fatalf("decoded wrong nums")
			}
		})
	}
}

func TestReadAllDeltaFast(t *testing.T) {
	for i := 0; i < 6; i++ {
		count := int(util.RandUint32() % 1e6)
		nums := util.GenUint32(count)
		util.SortUint32(nums)
		diffed := make([]uint32, count)
		util.Delta(nums, diffed)

		stream := writer.WriteAllScalar(diffed)
		t.Run(fmt.Sprintf("ReadAll: %d", count), func(t *testing.T) {
			out := make([]uint32, count)
			ReadAllDeltaFast(count, stream, out, 0)
			if !reflect.DeepEqual(nums, out) {
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
		stream := writer.WriteAllScalar(nums)
		out := make([]uint32, count)
		b.Run(fmt.Sprintf("Count_1e%d", i), func(b *testing.B) {
			b.SetBytes(int64(count * encode.MaxBytesPerNum))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				ReadAllFast(count, stream, out)
			}
			readSinkA = out
		})
	}
}

var readSinkB []uint32

func BenchmarkReadAllDeltaFast(b *testing.B) {
	for i := 0; i < 8; i++ {
		count := int(math.Pow10(i))
		nums := util.GenUint32(count)
		util.SortUint32(nums)
		stream := writer.WriteAllDeltaScalar(nums, 0)
		out := make([]uint32, count)
		b.Run(fmt.Sprintf("Count_1e%d", i), func(b *testing.B) {
			b.SetBytes(int64(count * encode.MaxBytesPerNum))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				ReadAllDeltaFast(count, stream, out, 0)
			}
			readSinkB = out
		})
	}
}

var readSinkC []uint32

func BenchmarkReadAllScalar(b *testing.B) {
	for i := 0; i < 8; i++ {
		count := int(math.Pow10(i))
		nums := util.GenUint32(count)
		util.SortUint32(nums)
		stream := writer.WriteAllScalar(nums)
		out := make([]uint32, count)
		b.Run(fmt.Sprintf("Count_1e%d", i), func(b *testing.B) {
			b.SetBytes(int64(count * encode.MaxBytesPerNum))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				ReadAllScalar(count, stream, out)
			}
			readSinkC = out
		})
	}
}

var readSinkD []uint32

func BenchmarkReadAllDeltaScalar(b *testing.B) {
	for i := 0; i < 8; i++ {
		count := int(math.Pow10(i))
		nums := util.GenUint32(count)
		util.SortUint32(nums)
		stream := writer.WriteAllDeltaScalar(nums, 0)
		out := make([]uint32, count)
		b.Run(fmt.Sprintf("Count_1e%d", i), func(b *testing.B) {
			b.SetBytes(int64(count * encode.MaxBytesPerNum))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				ReadAllDeltaScalar(count, stream, out, 0)
			}
			readSinkD = out
		})
	}
}

var readSinkE []uint32

func BenchmarkReadAllVarint(b *testing.B) {
	for i := 0; i < 8; i++ {
		count := int(math.Pow10(i))
		out := make([]uint32, count)
		data := make([]byte, binary.MaxVarintLen32*count)
		nums := util.GenUint32(count)
		util.SortUint32(nums)
		written := util.PutVarint(nums, data)
		data = data[:written]
		b.Run(fmt.Sprintf("Count_1e%d", i), func(b *testing.B) {
			b.SetBytes(int64(count * encode.MaxBytesPerNum))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				util.GetVarint(data, out)
			}
			readSinkB = out
		})
	}
}

var readSinkF []uint32

func BenchmarkReadAllDeltaVarint(b *testing.B) {
	for i := 0; i < 8; i++ {
		count := int(math.Pow10(i))
		out := make([]uint32, count)
		data := make([]byte, binary.MaxVarintLen32*count)
		nums := util.GenUint32(count)
		util.SortUint32(nums)
		written := util.PutDeltaVarint(nums, data, 0)
		data = data[:written]
		b.Run(fmt.Sprintf("Count_1e%d", i), func(b *testing.B) {
			b.SetBytes(int64(count * encode.MaxBytesPerNum))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				util.GetDeltaVarint(data, out, 0)
			}
			readSinkB = out
		})
	}
}
