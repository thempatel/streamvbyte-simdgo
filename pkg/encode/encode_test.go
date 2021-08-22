package encode

import (
	"reflect"
	"testing"
)

func TestPut8uint32Scalar(t *testing.T) {
	in := []uint32{1024, 3, 2, 1, 1_073_741_824, 10, 12, 1024}
	expectedData := []byte{
		0x00, 0x04, 0x03, 0x02, 0x01, 0x00, 0x00, 0x00, 0x40,
		0x0a, 0x0c, 0x00, 0x04,
	}
	
	expectedCtrl := uint16(0b01_00_00_11_00_00_00_01)
	out := make([]byte, 32)
	actualCtrl := put8uint32Scalar(in, out)
	if actualCtrl != expectedCtrl {
		t.Fatalf("expected: %#016b, got %#016b", expectedCtrl, actualCtrl)
	}

	actualData := out[:13]
	if !reflect.DeepEqual(expectedData, actualData) {
		t.Fatalf("expected %+v, got %+v", expectedData, actualData)
	}
}
