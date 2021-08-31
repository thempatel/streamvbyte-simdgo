package decode

import (
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

func TestGet8uint32Scalar(t *testing.T) {
	count := 8
	expected := util.GenUint32(count)
	in := make([]byte, count*encode.MaxBytesPerNum)
	ctrl := encode.Put8uint32Scalar(expected, in)
	out := make([]uint32, 8)

	read := Get8uint32Scalar(in, out, ctrl)
	if read != shared.ControlByteToSizeTwo(ctrl) {
		t.Fatalf("expected 13, got %d", read)
	}

	if !reflect.DeepEqual(expected, out) {
		t.Fatalf("expected %+v, got %+v", expected, out)
	}
}

func TestGet8uint32DiffScalar(t *testing.T) {
	count := 8
	expected := util.GenUint32(count)
	util.SortUint32(expected)
	in := make([]byte, count*encode.MaxBytesPerNum)
	ctrl := encode.Put8uint32DiffScalar(expected, in, 0)
	out := make([]uint32, 8)

	read := Get8uint32DiffScalar(in, out, ctrl, 0)
	if read != shared.ControlByteToSizeTwo(ctrl) {
		t.Fatalf("expected 13, got %d", read)
	}

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

	read := get8uint32(in, out, ctrl)
	if read != shared.ControlByteToSizeTwo(ctrl) {
		t.Fatalf("expected 13, got %d", read)
	}

	if !reflect.DeepEqual(expected, out) {
		t.Fatalf("expected %+v, got %+v", expected, out)
	}
}

func TestGet8uint32DiffFast(t *testing.T) {
	count := 8
	expected := util.GenUint32(count)
	util.SortUint32(expected)
	in := make([]byte, count*encode.MaxBytesPerNum)
	ctrl := encode.Put8uint32DiffScalar(expected, in, 0)
	out := make([]uint32, 8)

	read := get8uint32Diff(in, out, ctrl, 0)
	if read != shared.ControlByteToSizeTwo(ctrl) {
		t.Fatalf("expected 13, got %d", read)
	}

	if !reflect.DeepEqual(expected, out) {
		t.Fatalf("expected %+v, got %+v", expected, out)
	}
}

var readSinkA int

func BenchmarkGet8uint32Fast(b *testing.B) {
	count := 8
	nums := util.GenUint32(count)
	in := make([]byte, count*encode.MaxBytesPerNum)
	ctrl := encode.Put8uint32Scalar(nums, in)
	out := make([]uint32, count)

	read := 0
	b.SetBytes(32)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		read = get8uint32(in, out, ctrl)
	}
	readSinkA = read
}

var readSinkB int

func BenchmarkGet8uint32Scalar(b *testing.B) {
	count := 8
	nums := util.GenUint32(count)
	in := make([]byte, count*encode.MaxBytesPerNum)
	ctrl := encode.Put8uint32Scalar(nums, in)
	out := make([]uint32, count)

	read := 0
	b.SetBytes(32)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		read = Get8uint32Scalar(in, out, ctrl)
	}
	readSinkB = read
}

var readSinkC int

func BenchmarkGet8uint32Varint(b *testing.B) {
	data := util.PutVarint(util.GenUint32(8))

	read := 0
	b.SetBytes(32)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		read = util.GetVarint(data)
	}
	readSinkC = read
}
