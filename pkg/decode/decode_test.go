package decode

import (
	"reflect"
	"testing"
)

func TestGet8uint32Scalar(t *testing.T) {
	in := []byte{
		0x00, 0x04, 0x03, 0x02, 0x01, 0x00, 0x00, 0x00, 0x40,
		0x0a, 0x0c, 0x00, 0x04,
	}
	ctrl := uint16(0b01_00_00_11_00_00_00_01)
	expected := []uint32{1024, 3, 2, 1, 1_073_741_824, 10, 12, 1024}
	out := make([]uint32, 8)

	read := Get8uint32Scalar(in, out, ctrl)
	if read != 13 {
		t.Fatalf("expected 13, got %d", read)
	}

	if !reflect.DeepEqual(expected, out) {
		t.Fatalf("expected %+v, got %+v", expected, out)
	}
}

func TestGet8uint32Fast(t *testing.T) {
	in := []byte{
		0x00, 0x04, 0x03, 0x02, 0x01, 0x00, 0x00, 0x00, 0x40,
		0x0a, 0x0c, 0x00, 0x04,
	}
	for len(in) < 16 {
		in = append(in, 0x00)
	}

	ctrl := uint16(0b01_00_00_11_00_00_00_01)
	expected := []uint32{1024, 3, 2, 1, 1_073_741_824, 10, 12, 1024}
	out := make([]uint32, 8)

	read := get8uint32(in, out, ctrl)
	if read != 13 {
		t.Fatalf("expected 13, got %d", read)
	}

	if !reflect.DeepEqual(expected, out) {
		t.Fatalf("expected %+v, got %+v", expected, out)
	}
}
