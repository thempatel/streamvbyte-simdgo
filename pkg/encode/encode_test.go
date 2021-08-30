package encode

import (
	"encoding/binary"
	"reflect"
	"testing"

	"github.com/theMPatel/streamvbyte-simdgo/pkg/randutils"
	"github.com/theMPatel/streamvbyte-simdgo/pkg/shared"
)

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
		t.Fatalf("expected: %#016b, got %#016b", expectedCtrl, actualCtrl)
	}

	actualData := out[:13]
	if !reflect.DeepEqual(expectedData, actualData) {
		t.Fatalf("expected %+v, got %+v", expectedData, actualData)
	}
}

func TestPut8uint32Fast(t *testing.T) {
	if GetMode() == shared.Normal {
		t.Skipf("Testing environment doesn't support this test")
	}

	count := 8
	nums := make([]uint32, count)
	for i := 0; i < count; i++ {
		nums[i] = randutils.RandUint32()
	}

	out := make([]byte, MaxBytesPerNum*count)
	scalarCtrl := Put8uint32Scalar(nums, out)
	out = out[:shared.ControlByteToSizeTwo(scalarCtrl)]

	fastOut := make([]byte, MaxBytesPerNum*count)
	fastCtrl := put8uint32(nums, fastOut)
	fastOut = fastOut[:shared.ControlByteToSizeTwo(fastCtrl)]

	if scalarCtrl != fastCtrl {
		t.Fatalf("expected %#04x, actual %#04x", scalarCtrl, fastCtrl)
	}

	if !reflect.DeepEqual(out, fastOut) {
		t.Fatalf("expected %+v, got %+v", out, fastOut)
	}
}

var writeSinkA uint16

func BenchmarkPut8uint32Fast(b *testing.B) {
	count := 8
	nums := make([]uint32, count)
	for i := 0; i < count; i++ {
		nums[i] = randutils.RandUint32()
	}
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

func BenchmarkPut8uint32Scalar(b *testing.B) {
	count := 8
	nums := make([]uint32, count)
	for i := 0; i < count; i++ {
		nums[i] = randutils.RandUint32()
	}
	out := make([]byte, MaxBytesPerNum*count)

	var ctrl uint16
	b.SetBytes(32)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctrl = Put8uint32Scalar(nums, out)
	}
	writeSinkB = ctrl
}

var writeSinkC int

func BenchmarkPut8uint32Varint(b *testing.B) {
	count := 8
	nums := make([]uint32, count)
	for i := 0; i < count; i++ {
		nums[i] = randutils.RandUint32()
	}

	out := make([]byte, binary.MaxVarintLen32*count)
	written := 0

	b.SetBytes(32)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		written = write8Varint(nums, out)
	}
	writeSinkC = written
}

func write8Varint(nums []uint32, out []byte) int {
	pos := 0
	for i := range nums {
		size := binary.PutUvarint(out[pos:], uint64(nums[i]))
		pos += size
	}

	return pos
}