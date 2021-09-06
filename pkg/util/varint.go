package util

import "encoding/binary"

func PutVarint(nums []uint32, out []byte) int {
	pos := 0
	for i := range nums {
		size := binary.PutUvarint(out[pos:], uint64(nums[i]))
		pos += size
	}

	return pos
}

func GetVarint(data []byte, out []uint32) int {
	pos := 0
	i := 0
	for pos < len(data) {
		num, read := binary.Uvarint(data[pos:])
		pos += read
		out[i] = uint32(num)
		i++
	}
	return pos
}

func PutDeltaVarint(nums []uint32, out []byte, prev uint32) int {
	pos := 0
	for i := range nums {
		size := binary.PutUvarint(out[pos:], uint64(nums[i]-prev))
		pos += size
		prev = nums[i]
	}

	return pos
}

func GetDeltaVarint(in []byte, out []uint32, prev uint32) int {
	pos := 0
	i := 0
	for pos < len(in) {
		num, size := binary.Uvarint(in[pos:])
		pos += size
		res := uint32(num) + prev
		out[i] = res
		prev = res
		i++
	}

	return pos
}
