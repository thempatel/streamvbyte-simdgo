package x86

import (
	"fmt"
	"testing"
)

func TestWhatHappens(t *testing.T) {
	var nums []uint32
	for i := uint32(0); i < 8; i++ {
		nums = append(nums, i)
	}

	controlByte := x86ControlBytes8(nums)
	fmt.Printf("%b", controlByte)
}
