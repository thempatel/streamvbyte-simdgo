package util

import (
	"encoding/binary"
	"io"
)

func SilentClose(closer io.Closer) {
	_ = closer.Close()
}


func GenUint32(n int) []uint32 {
	nums := make([]uint32, n)
	for i := 0; i < n; i++ {
		nums[i] = RandUint32()
	}

	return nums
}

func PutVarint(nums []uint32) []byte {
	var (
		data []byte
		written int
		buf = make([]byte, binary.MaxVarintLen32)
	)

	for i := range nums {
		size := binary.PutUvarint(buf, uint64(nums[i]))
		data = append(data, buf[:size]...)
		written += size
	}

	return data[:written]
}

func GetVarint(data []byte) int {
	pos := 0
	for pos < len(data) {
		_, read := binary.Uvarint(data[pos:])
		pos += read
	}
	return pos
}