package decode

import (
	"encoding/binary"
	"math/rand"
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

var readSinkA int

func BenchmarkGet8uint32Fast(b *testing.B) {
	in := []byte{
		0x00, 0x04, 0x03, 0x02, 0x01, 0x00, 0x00, 0x00, 0x40,
		0x0a, 0x0c, 0x00, 0x04,
	}
	for len(in) < 16 {
		in = append(in, 0x00)
	}

	ctrl := uint16(0b01_00_00_11_00_00_00_01)
	out := make([]uint32, 8)

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
	in := []byte{
		0x00, 0x04, 0x03, 0x02, 0x01, 0x00, 0x00, 0x00, 0x40,
		0x0a, 0x0c, 0x00, 0x04,
	}

	ctrl := uint16(0b01_00_00_11_00_00_00_01)
	out := make([]uint32, 8)

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
	data := generate8Varint()

	read := 0
	b.SetBytes(32)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		read = decode8Varint(data)
	}
	readSinkC = read
}

func generate8Varint() []byte {
	count := 8
	var (
		data []byte
		written int
		buf = make([]byte, binary.MaxVarintLen32)
	)

	for i := 0; i < count; i++ {
		size := binary.PutUvarint(buf, uint64(rand.Uint32()))
		data = append(data, buf[:size]...)
		written += size
	}

	return data[:written]
}

func decode8Varint(data []byte) int {
	pos := 0
	for i := 0; i < 8; i++ {
		_, read := binary.Uvarint(data[pos:])
		pos += read
	}
	return pos
}
