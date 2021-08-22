// +build amd64

package decode

import "golang.org/x/sys/cpu"

var check cpuCheck = func() bool {
	return cpu.X86.HasAVX
}

func get8uint32(in []byte, out []uint32, ctrl uint8) int {
	return 0
}