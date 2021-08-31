package util

import (
	"encoding/binary"
	"io"
	"sort"
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

func PutDiffVarint(nums []uint32, out []byte, prev uint32) int {
	pos := 0
	for i := range nums {
		size := binary.PutUvarint(out[pos:], uint64(nums[i]-prev))
		pos += size
		prev = nums[i]
	}

	return pos
}

func GetDiffVarint(in []byte, out []uint32, prev uint32) int {
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

func SortUint32(in []uint32) {
	sort.Slice(in, func(i, j int) bool {
		return in[i] < in[j]
	})
}

func Diff(in []uint32, out []uint32) {
	for i := range in {
		if i > 0 {
			out[i] = in[i] - in[i-1]
		} else {
			out[i] = in[i]
		}
	}
}
