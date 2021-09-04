package decode

import (
	"encoding/binary"
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

func TestGet8uint32Scalar(t *testing.T) {
	count := 8
	expected := util.GenUint32(count)
	in := make([]byte, count*encode.MaxBytesPerNum)
	ctrl := encode.Put8uint32Scalar(expected, in)
	out := make([]uint32, 8)

	Get8uint32Scalar(in, out, ctrl)
	if !reflect.DeepEqual(expected, out) {
		t.Fatalf("expected %+v, got %+v", expected, out)
	}
}

func TestGet8uint32DeltaScalar(t *testing.T) {
	count := 8
	expected := util.GenUint32(count)
	util.SortUint32(expected)
	in := make([]byte, count*encode.MaxBytesPerNum)
	ctrl := encode.Put8uint32DeltaScalar(expected, in, 0)
	out := make([]uint32, 8)

	Get8uint32DeltaScalar(in, out, ctrl, 0)
	if !reflect.DeepEqual(expected, out) {
		t.Fatalf("expected %+v, got %+v", expected, out)
	}
}

func TestGet8uint32Fast(t *testing.T) {
	count := 8
	expected := util.GenUint32(count)
	in := make([]byte, count*encode.MaxBytesPerNum)
	ctrl := encode.Put8uint32Scalar(expected, in)
	out := make([]uint32, 8)

	Get8uint32Fast(in, out, ctrl)
	if !reflect.DeepEqual(expected, out) {
		t.Fatalf("expected %+v, got %+v", expected, out)
	}
}

func TestGet8uint32DeltaFast(t *testing.T) {
	count := 8
	expected := util.GenUint32(count)
	util.SortUint32(expected)
	in := make([]byte, count*encode.MaxBytesPerNum)
	ctrl := encode.Put8uint32DeltaScalar(expected, in, 0)
	out := make([]uint32, 8)

	Get8uint32DeltaFast(in, out, ctrl, 0)
	if !reflect.DeepEqual(expected, out) {
		t.Fatalf("expected %+v, got %+v", expected, out)
	}
}

var readSinkA []uint32

func BenchmarkGet8uint32Fast(b *testing.B) {
	count := 8
	out := make([]uint32, count)

	nums := util.GenUint32(count)
	in := make([]byte, count*encode.MaxBytesPerNum)
	ctrl := encode.Put8uint32Scalar(nums, in)

	b.SetBytes(int64(count*encode.MaxBytesPerNum))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Get8uint32Fast(in, out, ctrl)
	}
	readSinkA = out
}

var readSinkB []uint32

func BenchmarkGet8uint32DeltaFast(b *testing.B) {
	count := 8
	out := make([]uint32, count)
	nums := util.GenUint32(count)
	util.SortUint32(nums)
	in := make([]byte, count*encode.MaxBytesPerNum)
	ctrl := encode.Put8uint32DeltaScalar(nums, in, 0)

	b.SetBytes(int64(count*encode.MaxBytesPerNum))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Get8uint32DeltaFast(in, out, ctrl, 0)
	}
	readSinkB = out
}

var readSinkC []uint32

func BenchmarkGet8uint32Scalar(b *testing.B) {
	count := 8
	out := make([]uint32, count)
	nums := util.GenUint32(count)
	in := make([]byte, count*encode.MaxBytesPerNum)
	ctrl := encode.Put8uint32Scalar(nums, in)

	b.SetBytes(int64(count*encode.MaxBytesPerNum))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Get8uint32Scalar(in, out, ctrl)
	}
	readSinkC = out
}

var readSinkD []uint32

func BenchmarkGet8uint32DeltaScalar(b *testing.B) {
	count := 8
	out := make([]uint32, count)
	nums := util.GenUint32(count)
	util.SortUint32(nums)
	in := make([]byte, count*encode.MaxBytesPerNum)
	ctrl := encode.Put8uint32DeltaScalar(nums, in, 0)

	b.SetBytes(int64(count*encode.MaxBytesPerNum))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Get8uint32DeltaScalar(in, out, ctrl, 0)
	}
	readSinkD = out
}

var readSinkE []uint32

func BenchmarkGet8uint32Varint(b *testing.B) {
	count := 8
	out := make([]uint32, count)
	data := make([]byte, binary.MaxVarintLen32*count)
	written := util.PutVarint(util.GenUint32(count), data)
	data = data[:written]

	b.SetBytes(int64(count*encode.MaxBytesPerNum))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		util.GetVarint(data, out)
	}
	readSinkE = out
}

var readSinkF []uint32

func BenchmarkGet8uint32DeltaVarint(b *testing.B) {
	count := 8
	out := make([]uint32, count)
	data := make([]byte, binary.MaxVarintLen32*count)
	written := util.PutDeltaVarint(util.GenUint32(count), data, 0)
	data = data[:written]

	b.SetBytes(int64(count*encode.MaxBytesPerNum))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		util.GetDeltaVarint(data, out, 0)
	}
	readSinkF = out
}