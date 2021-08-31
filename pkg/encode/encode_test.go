package encode

import (
	"encoding/binary"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/theMPatel/streamvbyte-simdgo/pkg/shared"
	"github.com/theMPatel/streamvbyte-simdgo/pkg/util"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestPut8uint32Scalar(t *testing.T) {
	in := []uint32{1024, 3, 2, 1, 1_073_741_824, 10, 12, 1024}
	expectedData := []byte{
		0x00, 0x04, 0x03, 0x02, 0x01, 0x00, 0x00, 0x00, 0x40,
		0x0a, 0x0c, 0x00, 0x04,
	}

	expectedCtrl := uint16(0b01_00_00_11_00_00_00_01)
	out := make([]byte, 32)
	actualCtrl := Put8uint32Scalar(in, out)
	if actualCtrl != expectedCtrl {
		t.Fatalf("expected: %#016b, got %#016b, %+v", expectedCtrl, actualCtrl, in)
	}

	actualData := out[:13]
	if !reflect.DeepEqual(expectedData, actualData) {
		t.Fatalf("expected %+v, got %+v, %+v", expectedData, actualData, in)
	}
}

func TestPut8uint32DiffScalar(t *testing.T) {
	count := 8
	nums := util.GenUint32(count)
	util.SortUint32(nums)
	diffed := make([]uint32, count)
	util.Diff(nums, diffed)

	expectedData := make([]byte, count*MaxBytesPerNum)
	expectedCtrl := Put8uint32Scalar(diffed, expectedData)
	expectedData = expectedData[:shared.ControlByteToSizeTwo(expectedCtrl)]

	out := make([]byte, count*MaxBytesPerNum)
	actualCtrl := Put8uint32DiffScalar(nums, out, 0)
	if actualCtrl != expectedCtrl {
		t.Fatalf("expected: %#016b, got %#016b, %+v", expectedCtrl, actualCtrl, nums)
	}

	actualData := out[:shared.ControlByteToSizeTwo(actualCtrl)]
	if !reflect.DeepEqual(expectedData, actualData) {
		t.Fatalf("expected %+v, got %+v, %+v", expectedData, actualData, nums)
	}
}

func TestPut8uint32Fast(t *testing.T) {
	if GetMode() == shared.Normal {
		t.Skipf("Testing environment doesn't support this test")
	}

	count := 8
	nums := util.GenUint32(count)

	out := make([]byte, MaxBytesPerNum*count)
	scalarCtrl := Put8uint32Scalar(nums, out)
	out = out[:shared.ControlByteToSizeTwo(scalarCtrl)]

	fastOut := make([]byte, MaxBytesPerNum*count)
	fastCtrl := put8uint32(nums, fastOut)
	fastOut = fastOut[:shared.ControlByteToSizeTwo(fastCtrl)]

	if scalarCtrl != fastCtrl {
		t.Fatalf("expected %#04x, actual %#04x, %+v", scalarCtrl, fastCtrl, nums)
	}

	if !reflect.DeepEqual(out, fastOut) {
		t.Fatalf("expected %+v, got %+v, %+v", out, fastOut, nums)
	}
}

func TestPut8uint32DiffFast(t *testing.T) {
	if GetMode() == shared.Normal {
		t.Skipf("Testing environment doesn't support this test")
	}

	count := 8
	nums := util.GenUint32(count)
	util.SortUint32(nums)

	expectedData := make([]byte, MaxBytesPerNum*count)
	scalarCtrl := Put8uint32DiffScalar(nums, expectedData, 0)
	expectedData = expectedData[:shared.ControlByteToSizeTwo(scalarCtrl)]

	fastOut := make([]byte, MaxBytesPerNum*count)
	fastCtrl := put8uint32Diff(nums, fastOut, 0)
	fastOut = fastOut[:shared.ControlByteToSizeTwo(fastCtrl)]

	if scalarCtrl != fastCtrl {
		t.Fatalf("expected %#04x, actual %#04x, %+v", scalarCtrl, fastCtrl, nums)
	}

	if !reflect.DeepEqual(expectedData, fastOut) {
		t.Fatalf("expected %+v, got %+v, %+v", expectedData, fastOut, nums)
	}
}

var writeSinkA uint16

func BenchmarkPut8uint32Fast(b *testing.B) {
	count := 8
	nums := util.GenUint32(count)
	out := make([]byte, MaxBytesPerNum*count)

	var ctrl uint16
	b.SetBytes(32)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctrl = put8uint32(nums, out)
	}
	writeSinkA = ctrl
}

var writeSinkB uint16

func BenchmarkPut8uint32DiffFast(b *testing.B) {
	count := 8
	nums := util.GenUint32(count)
	util.SortUint32(nums)
	out := make([]byte, MaxBytesPerNum*count)

	var ctrl uint16
	b.SetBytes(32)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctrl = put8uint32Diff(nums, out, 0)
	}
	writeSinkB = ctrl
}

var writeSinkC uint16

func BenchmarkPut8uint32Scalar(b *testing.B) {
	count := 8
	nums := util.GenUint32(count)
	out := make([]byte, MaxBytesPerNum*count)

	var ctrl uint16
	b.SetBytes(32)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctrl = Put8uint32Scalar(nums, out)
	}
	writeSinkC = ctrl
}

var writeSinkD uint16

func BenchmarkPut8uint32DiffScalar(b *testing.B) {
	count := 8
	nums := util.GenUint32(count)
	util.SortUint32(nums)
	out := make([]byte, MaxBytesPerNum*count)

	var ctrl uint16
	b.SetBytes(32)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctrl = Put8uint32DiffScalar(nums, out, 0)
	}
	writeSinkD = ctrl
}

var writeSinkE int

func BenchmarkPut8uint32Varint(b *testing.B) {
	count := 8
	nums := util.GenUint32(count)

	out := make([]byte, binary.MaxVarintLen32*count)
	written := 0

	b.SetBytes(32)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		written = write8Varint(nums, out)
	}
	writeSinkE = written
}

var writeSinkF int

func BenchmarkPut8uint32DiffVarint(b *testing.B) {
	count := 8
	nums := util.GenUint32(count)
	util.SortUint32(nums)
	diffed := make([]uint32, count)
	util.Diff(nums, diffed)

	out := make([]byte, binary.MaxVarintLen32*count)
	written := 0

	b.SetBytes(32)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		written = write8Varint(diffed, out)
	}
	writeSinkF = written
}

func write8Varint(nums []uint32, out []byte) int {
	pos := 0
	for i := range nums {
		size := binary.PutUvarint(out[pos:], uint64(nums[i]))
		pos += size
	}

	return pos
}
