// +build amd64

// Package encode provides an x86_64 implementation of two
// Stream VByte encoding algorithms, a normal encoding approach
// and one that incorporates differential coding.
package encode

import (
	"github.com/theMPatel/streamvbyte-simdgo/pkg/shared"
	"golang.org/x/sys/cpu"
)

// GetMode performs a check to see if the current ISA supports
// the below encoding funcs.
func GetMode() shared.PerformanceMode {
	if cpu.X86.HasAVX && cpu.X86.HasAVX2 {
		return shared.Fast
	}
	return shared.Normal
}

// put8uint32 binds to put8uint32Fast which is implemented
// in assembly.
func put8uint32(in []uint32, out []byte) uint16 {
	return put8uint32Fast(in, out,
		shared.EncodeShuffleTable,
		shared.PerControlLenTable,
	)
}

// put8uint32Diff binds to put8uint32DiffFast which is implemented
// in assembly.
func put8uint32Diff(in []uint32, out []byte, prev uint32) uint16 {
	return put8uint32DiffFast(
		in, out, prev,
		shared.EncodeShuffleTable,
		shared.PerControlLenTable,
	)
}

// put8uint32Fast has three core phases. First a 16-bit control is
// generated for the incoming 8 uint32s. Then, the calculated control
// is used to index into shared.EncodeShuffleTable to fetch the
// correct shuffle mask to compress the incoming integers. Finally,
// the calculated control is used to index into shared.PerControlLenTable
// to determine the offsets in the output array to write to.
//
// Based on the algorithm devised by Lemire et al., the SIMD control
// byte generation algorithm proceeds as follows. Note that here we
// are using 1234 as our first example integer.
//
// 00000000 00000000 00000100 11010010 // 1234
// 00000001 00000001 00000001 00000001 // 0x0101 mask
// ----------------------------------- // byte-min(1234, 0x0101)
// 00000000 00000000 00000001 00000001
//
// The algorithm first uses a mask where every byte is equal to 1. If
// you perform a per-byte min operation on our integer and the 1's mask,
// the result will have a 1 at every byte that had a value in the original
// integer.
//
// 00000000 00000000 00000001 00000001
// ----------------------------------- // pack unsigned saturating
// 00000000 00000000 00000000 11111111 // 16-bit to 8-bit
//
// Now you perform a 16-bit to 8-bit unsigned saturating pack operation.
// Practically this means that you're taking every 16-bit value and trying
// to shove that into 8 bits. If the 16-bit integer is larger than the
// largest unsigned integer 8 bits can support, the pack saturates to the
// largest unsigned 8-bit value.
//
// Why this is performed will become more clear in the subsequent steps,
// however, at a high level, for every integer you want to encode, you
// want for the MSB of two consecutive bytes in the control bits stream
// to be representative of the final 2-bit control. For example, if you
// have a 3-byte integer, you want the MSB of two consecutive bytes to be
// 1 and 0, in that order. The reason you would want this is that there
// is a vector pack instruction that takes the MSB from every byte in the
// control bits stream and packs it into the lowest byte. This would thus
// represent the value `0b10` in the final byte for this 3-byte integer,
// which is what we want.
//
// Performing a 16-bit to 8-bit unsigned saturating pack has the effect
// that you can use the saturation behavior to conditionally turn on the
// MSB of these bytes depending on which bytes have values in the original
// 32-bit integer.
//
// 00000000 00000000 00000000 11111111 // control bits stream
// 00000001 00000001 00000001 00000001 // 0x0101 mask
// ----------------------------------- // signed 16-bit min
// 00000000 00000000 00000000 11111111
//
// We then take the 1's mask we used before and perform a signed 16-bit
// min operation. The reason for this is more clear if you look at an
// example using a 3-byte integer.
//
// 00000000 00001100 00001010 10000011 // 789123
// 00000001 00000001 00000001 00000001 // 0x0101 mask
// ----------------------------------- // byte-min(789123, 0x0101)
// 00000000 00000001 00000001 00000001
// ----------------------------------- // pack unsigned saturating 16-bit to 8-bit
// 00000000 00000000 00000001 11111111
// 00000001 00000001 00000001 00000001 // 0x0101 mask
// ----------------------------------- // signed 16-bit min
// 00000000 00000000 00000001 00000001
//
// The signed 16-bit min operation has three important effects.
//
// First, for 3-byte integers, it has the effect of turning off the
// MSB of the lowest byte. This is necessary because a 3-byte integer
// should have a 2-bit control that is `0b10` and without this step
// using the MSB pack operation would result in a 2-bit control that
// looks something like `0b_1`, where the lowest bit is on. Obviously
// this is wrong, since only integers that require 2 or 4 bytes to
// encode should have that lower bit on, i.e. 1 or 3 as a zero-indexed
// length.
//
// Second, for 4-byte integers, the signed aspect has the effect of
// leaving both MSBs of the 2 bytes on. When using the MSB pack
// operation later on, it will result in a 2-bit control value of
// `0b11`, which is what we want.
//
// Third, for 1 and 2 byte integers, it has no effect. This is great
// for 2-byte values since the MSB will remain on and 1 byte values
// will not have any MSB on anyways, so it is effectively a noop in
// both scenarios.
//
// 00000000 00000000 00000000 11111111 // control bits stream (original 1234)
// 01111111 00000000 01111111 00000000 // 0x7F00 mask
// ----------------------------------- // add unsigned saturating 16-bit
// 01111111 00000000 01111111 11111111
//
// Next, we take a mask with the value `0x7F00` and perform an unsigned
// saturating add to the control bits stream. In the case for the integer
// `1234` this has no real effect. We maintain the MSB in the lowest byte.
// You'll note, however, that the only byte that has its MSB on is the last
// one, so performing an MSB pack operation would result in a value of
// `0b0001`, which is what we want. An example of this step on the integer
// `789123` might paint a clearer picture.
//
// 00000000 00000000 00000001 00000001 // control bits stream (789123)
// 01111111 00000000 01111111 00000000 // 0x7F00 mask
// ----------------------------------- // add unsigned saturating 16-bit
// 01111111 00000000 11111111 00000001
//
// You'll note here that the addition of `0x01` with `0x7F` in the upper
// byte results in the MSB of the resulting upper byte turning on. The MSB
// in the lower byte remains off and now an MSB pack operation will resolve
// to `0b0010`, which is what we want. The unsigned saturation behavior is
// really important for 4-byte numbers that only have bits in the most
// significant byte on. An example below:
//
// 01000000 00000000 00000000 00000000 // 1073741824
// 00000001 00000001 00000001 00000001 // 0x0101 mask
// ----------------------------------- // byte-min(1073741824, 0x0101)
// 00000001 00000000 00000000 00000000
// ----------------------------------- // pack unsigned saturating 16-bit to 8-bit
// 00000000 00000000 11111111 00000000
// 00000001 00000001 00000001 00000001 // 0x0101 mask
// ----------------------------------- // signed 16-bit min
// 00000000 00000000 11111111 00000000
// 01111111 00000000 01111111 00000000 // 0x7F00 mask
// ----------------------------------- // add unsigned saturating 16-bit
// 01111111 00000000 11111111 11111111
//
// Note here that because only the upper byte had a value in it, the lowest
// byte in the control bits stream remains zero for the duration of the
// algorithm. This poses an issue, since for a 4-byte value, we want for the
// 2-bit control to result in a value of `0b11`. Performing a 16-bit unsigned
// saturating addition has the effect of turning on all bits in the lower
// byte, and thus we get a result with the MSB in the lower byte on.
//
// 01111111 00000000 11111111 00000001 // control bits stream (789123)
// ----------------------------------- // move byte mask
// 00000000 00000000 00000000 00000010 // 2-bit control
//
// The final move byte mask is performed on the control bits stream, and we
// now have the result we wanted.
//
// We then use the above control bits we generated to get the appropriate
// shuffle masks. Below is an example of how the shuffle operation and a
// mask allows for us to tightly pack two integers into the output buffer.
//
// input [1234, 789123] (little endian R-to-L)
// 00000000 00001100 00001010 10000011 00000000 00000000 00000100 11010010
//            |       |         |                             |        |
//            |       |         |____________________         |        |
//            |       |_____________________         |        |        |
//            |____________________         |        |        |        |
//                                 v        v        v        v        v
//    0xff     0xff     0xff     0x06     0x05     0x04     0x01     0x00 // mask in hex
// -----------------------------------------------------------------------
// 00000000 00000000 00000000 00001100 00001010 10000011 00000100 11010010 // packed
//go:noescape
func put8uint32Fast(
	in []uint32, outBytes []byte,
	shuffle *[256][16]uint8, lenTable *[256]uint8,
) (r uint16)

// put8uint32DiffFast works similarly to put8uint32Fast above, except
// that prior to encoding the 8 uint32s, we first use differential
// coding to change the original numbers into deltas using SIMD
// techniques. Afterwards, the encoding algorithm follows the same
// flow as put8uint32Fast. The basic differential coding algorithm
// is as follows:
//
// Prev:  			[P P P P]
// Input: 			[A B C D]
// Concat-shift:	[P A B C]
// Subtract:		[A-P B-A C-B D-C]
//go:noescape
func put8uint32DiffFast(
	in []uint32, outBytes []byte, prev uint32,
	shuffle *[256][16]uint8, lenTable *[256]uint8,
) (r uint16)
