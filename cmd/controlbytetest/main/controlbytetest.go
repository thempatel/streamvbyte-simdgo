package main

import (
	"fmt"

	"github.com/theMPatel/streamvbyte-simdgo/internal/encode"
)

func main() {
	var nums []uint32
	for i := uint32(0); i < 8; i++ {
		nums = append(nums, i)
	}

	controlByte := encode.ControlBytes(nums)
	fmt.Printf("%b", controlByte)
}
