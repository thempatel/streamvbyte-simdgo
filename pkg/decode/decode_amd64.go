// +build amd64

// Package decode provides an x86_64 implementation of two
// Stream VByte decoding algorithms, a normal decoding approach
// and one that incorporates differential coding.
package decode

import (
	"github.com/theMPatel/streamvbyte-simdgo/pkg/shared"
	"golang.org/x/sys/cpu"
)

// GetMode performs a check to see if the current ISA supports
// the below decoding funcs.
func GetMode() shared.PerformanceMode {
	if cpu.X86.HasAVX {
		return shared.Fast
	}
	return shared.Normal
}

// Get8uint32Fast binds to get8uint32Fast which is implemented in
// assembly.
func Get8uint32Fast(in []byte, out []uint32, ctrl uint16) {
	get8uint32Fast(in, out, ctrl,
		shared.DecodeShuffleTable,
		shared.PerControlLenTable,
	)
}

// Get8uint32DeltaFast binds to get8uint32DeltaFast which is implemented
// in assembly.
func Get8uint32DeltaFast(in []byte, out []uint32, ctrl uint16, prev uint32) {
	get8uint32DeltaFast(
		in, out, ctrl, prev,
		shared.DecodeShuffleTable,
		shared.PerControlLenTable,
	)
}

// get8uint32Fast uses the provided 16-bit control to load the
// appropriate decoding shuffle masks and performs a shuffle
// operation on the provided input bytes. This in effect decompresses
// the input byte stream to uint32s. The result is written to
// the provided output slice.
//go:noescape
func get8uint32Fast(
	in []byte, out []uint32, ctrl uint16,
	shuffle *[256][16]uint8, lenTable *[256]uint8,
)

// get8uint32DeltaFast works similarly to get8uint32Fast with the
// exception that prior to writing the uncompressed integers out
// to the output slice, the original values are reconstructed from
// the diffs. The basic reconstruction algorithm is as follows:
//
// Input:           [A B C D]
// Input Shifted:   [- A  B  C]
// Add above two:   [A AB BC CD]
// Add Prev:        [PA PAB PBC PCD]
// Input Shifted:   [- - A AB]
// Add Shifted:     [PA PAB PABC PABCD]
//go:noescape
func get8uint32DeltaFast(
	in []byte, out []uint32, ctrl uint16, prev uint32,
	shuffle *[256][16]uint8, lenTable *[256]uint8,
)
