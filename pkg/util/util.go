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

func SortUint32(in []uint32) {
	sort.Slice(in, func(i, j int) bool {
		return in[i] < in[j]
	})
}

func Delta(in []uint32, out []uint32) {
	for i := range in {
		if i > 0 {
			out[i] = in[i] - in[i-1]
		} else {
			out[i] = in[i]
		}
	}
}

// MakeOutByteArr serves to create output arrays that occupy different parts
// of the processor's cache line. This is done to ensure that the
// benchmarks more accurately reflect real world performance characteristics.
// Reusing the same output array in the benchmarks has the effect of
// consistently hitting L1 cache, and to a lesser extent L2 and L3. Additionally,
// depending on the array size, there's a good chance writes all hit the same
// cache line which ultimately inflates the benchmark numbers
//
// An L1 cache size may likely be 64 KB and a cache line might be 64 bytes. Choose
// a count and total number with these numbers in mind
func MakeOutByteArr(count, total int) [][]byte {
	toRet := make([][]byte, total)
	for i := range toRet {
		toRet[i] = make([]byte, count)
	}

	return toRet
}

// MakeOutUint32Arr serves the same purpose as MakeOutByteArr.
func MakeOutUint32Arr(count, total int) [][]uint32 {
	toRet := make([][]uint32, total)
	for i := range toRet {
		toRet[i] = make([]uint32, count)
	}

	return toRet
}