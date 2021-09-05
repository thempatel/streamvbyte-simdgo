package pkg

import (
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/theMPatel/streamvbyte-simdgo/pkg/decode"
	"github.com/theMPatel/streamvbyte-simdgo/pkg/encode"
	"github.com/theMPatel/streamvbyte-simdgo/pkg/util"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestRoundTripScalar(t *testing.T) {
	in := []uint32{1024, 3, 2, 1, 1_073_741_824, 10, 12, 1024}
	expectedData := []byte{
		0x00, 0x04, 0x03, 0x02, 0x01, 0x00, 0x00, 0x00, 0x40,
		0x0a, 0x0c, 0x00, 0x04,
	}

	expectedCtrl := uint16(0b01_00_00_11_00_00_00_01)
	out := make([]byte, 32)
	actualCtrl := encode.Put8uint32Scalar(in, out)
	if actualCtrl != expectedCtrl {
		t.Fatalf("expected: %#016b, got %#016b", expectedCtrl, actualCtrl)
	}

	actualData := out[:13]
	if !reflect.DeepEqual(expectedData, actualData) {
		t.Fatalf("expected %+v, got %+v", expectedData, actualData)
	}

	decoded := make([]uint32, 8)
	decode.Get8uint32Scalar(actualData, decoded, actualCtrl)

	if !reflect.DeepEqual(in, decoded) {
		t.Fatalf("expected %+v, actual %+v", in, decoded)
	}
}

func BenchmarkMemCopy8Uint32(b *testing.B) {
	count := 8
	nums := make([]uint32, count)
	for i := 0; i < count; i++ {
		nums[i] = util.RandUint32()
	}

	out := make([]uint32, count)
	b.SetBytes(int64(count * encode.MaxBytesPerNum))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		copy(out, nums)
	}
}
