package encode

import (
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/theMPatel/streamvbyte-simdgo/pkg/shared"
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
	for i := 0; i < 8; i++ {
		nums[i] = rand.Uint32()
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

var ctrlSinkA uint16

func BenchmarkPut8uint32Fast(b *testing.B) {
	count := 8
	nums := make([]uint32, count)
	for i := 0; i < count; i++ {
		nums[i] = rand.Uint32()
	}
	out := make([]byte, MaxBytesPerNum*count)

	var ctrl uint16
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctrl = put8uint32(nums, out)
	}
	ctrlSinkA = ctrl
}

var ctrlSinkB uint16

func BenchmarkPut8uint32Scalar(b *testing.B) {
	count := 8
	nums := make([]uint32, count)
	for i := 0; i < 8; i++ {
		nums[i] = rand.Uint32()
	}
	out := make([]byte, MaxBytesPerNum*count)

	var ctrl uint16
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctrl = Put8uint32Scalar(nums, out)
	}
	ctrlSinkB = ctrl
}
