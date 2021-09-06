package util

import (
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
