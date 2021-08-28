# Stream VByte SIMD Go

This is a repository that contains a port of Stream VByte to Go. Notably, this repo takes extra care
to leverage SIMD techniques to achieve better performance. Currently, there is support for x86_64 architectures
that have AVX and AVX2 hardware support. In cases where that is not available, or on non x86_64 architectures
there is a portable scalar implementation. We also perform a runtime check to make sure that the necessary
ISA are available and if not fallback to the scalar approach.

There are several existing implementations:

1. [Reference Implementation](https://github.com/lemire/streamvbyte)
2. [Rust Implementation](https://bitbucket.org/marshallpierce/stream-vbyte-rust)

There is also another Go version [here](https://github.com/nelz9999/stream-vbyte-go) however, it only has a
scalar implementation which prompted this implementation with SIMD techniques.

---
Stream VByte uses the same underlying format as Google's Group Varint approach. Lemire et al. wanted
to see if there was away to improve the performance even more and introduced a clever twist to enable
better performance via SIMD techniques.

The basic goal of the Group Varint format is to be able to compress integers and load them 
really quickly. This has two advantages, you save space on disk and consequently save time
when loading the compressed data into memory, since there is less of it. The way they achieve
this compression is by improving upon the insight that backs a more basic Varint encoding.

## Varint format

The insight that backs the more basic Varint encoding is noticing that you oftentimes
don't need 32 bits to encode a 32-bit integer. Take for example an unsigned number that is less
than 2^8 (256). This number will have bits set in the lowest byte of a 32-bit number, while the
remaining 3 bytes will simply be zeros.

```go
package foo

// Num in binary:
//
// 00000000 00000000 00000000 01101111
var Num uint32 = 111
```

An approach you can take to compress this number is to encode the number using a variable
number of bytes. For example, you can use the lower 7 bits to store data, i.e. bits
from the original number, and then use the MSB as a continuation bit. If the MSB bit is on, i.e.
is 1, then more bytes are needed to decode this particular number. Below is an example where
you might need 2 bytes to store the number 1234.

```go
package foo

// Num in binary:
//
// 00000000 00000000 00000100 11010010
//
// Num compressed:
//
// 0|0001001| 1|1010010|
// ^          ^ Continuation bits
//
//     ^           ^ Data bits
var Num uint32 = 1234
```

If you want to decode this number, you simply build up the number iteratively. I.e. you OR the
last 7 bits of every byte shifted to the appropriate length to your 32-bit number until you
find a byte that doesn't have a continuation bit set. Note that this works the same for 64-bit
numbers.

The problem with this approach is that it introduces a lot of branch mis-predictions during encoding/decoding.
During the decoding phase, you don't know ahead of time the number of bytes that were used to encode the number
you are currently processing and so you need to iterate until you find a byte without a continuation bit on.
If you have numbers that are nonuniform, i.e. numbers that require random numbers of bytes to encode relative
to one another, this can pose a challenge to the processor's branch predictor. These mis-predictions can cause
major slowdowns in processor pipelines and so was born the Group Varint format.

## Group Varint format

The Group Varint format assumes that everything you hope to achieve, you can do with 32-bit numbers.
It introduces the concept of a control byte which is simply a byte that stores the encoded
lengths of a group of 4 32-bit numbers, hence Group Varint. 32-bit numbers only require up to 4 bytes
to properly encode. This means that you can represent their lengths with 2 bits using a zero-indexed length
i.e. 0, 1, 2, and 3 to represent numbers that require 1, 2, 3 and 4 bytes to encode, respectively.

```go
package foo

// Nums in binary:
//
// 00000000 00000000 00000000 01101111  =        111 
// 00000000 00000000 00000100 11010010  =       1234
// 00000000 00001100 00001010 10000011  =     789123
// 01000000 00000000 00000000 00000000  = 1073741824
//
// Num         Len      Control byte
// ---------------------------------
// 111          1               0b00 
// 1234         2               0b01
// 789123       3               0b10
// 1073741824   4               0b11
//
// Final Control byte
// 0b11100100
//
// Encoded data (little endian right-to-left) 
// 0b01000000 0b00000000 0b00000000 0b00000000 0b00001100 0b00001010 0b10000011 0b00000100 0b11010010 0b01101111
var (
	Num1 uint32 = 111
	Num2 uint32 = 1234
	Num3 uint32 = 789_123
	Num4 uint32 = 1_073_741_824
)
```

You can then prefix every group of 4 encoded 32-bit numbers with their control byte and then use it during decoding.
The obvious downside is that you pay a storage cost of one byte for every 4 numbers you want to encode. For 2^20 
encoded numbers, that's an extra 256 KB of extra space: totally marginal. The great upside, though, is that
you've now removed almost all branches from your decoding phase. You know exactly how many data bytes you need
to read from a buffer for a particular number and then can use branchless decoding.

```go
package foo

import (
	"encoding/binary"
)

func decodeOne(input []byte, size uint8) uint32 {
	buf := make([]byte, 4)
	copy(buf, input[:size])

	 // func (littleEndian) Uint32(b []byte) uint32 {
	 // 	_ = b[3]
	 // 	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
	 // }
	return binary.LittleEndian.Uint32(buf)
}

func main() {
	ctrl := uint8(0b11100100)
	data := []byte{
		0b01101111, 0b11010010, 0b00000100,
		0b10000011, 0b00001010, 0b00001100,
		0b00000000, 0b00000000, 0b00000000,
		0b01000000, 
    }
    
	len0 := (ctrl & 3) + 1      // 1
	len1 := (ctrl >> 2 & 3) + 1 // 2
	len2 := (ctrl >> 4 & 3) + 1 // 3
	len3 := (ctrl >> 6 & 3) + 1 // 4
	
	_ = decodeOne(data, len0) // 111
	_ = decodeOne(data, len1) // 1234
	_ = decodeOne(data, len2) // 789_123
	_ = decodeOne(data, len3) // 1_073_741_824
}
```

# References

[Stream VByte: Faster Byte-Oriented Integer Compression](https://arxiv.org/pdf/1709.08990.pdf)

// 2-byte
// 00000000 00000000 00000010 00001010 Value
// 00000001 00000001 00000001 00000001 1's Mask
// 00000000 00000000 00000001 00000001 min
// 00000000 00000000 00000000 11111111 pack
// 00000001 00000001 00000001 00000001 1's Mask
// 00000000 00000000 00000000 11111111 16bit signed min
// 01111111 00000000 01111111 00000000 7F00-mask
// 01111111 00000000 01111111 11111111 Add 7F00-mask
// 00000000 00000000 00000000 00000001 movemask



// 4-byte
// 00000001 00000000 00000000 00000000 Value
// 00000001 00000001 00000001 00000001 1's Mask
// 00000001 00000000 00000000 00000000 min
// 00000000 00000000 11111111 00000000 pack
// 00000001 00000001 00000001 00000001 1's Mask
// 00000000 00000000 11111111 00000000 16bit signed min
// 01111111 00000000 01111111 00000000 7F00-mask
// 01111111 00000000 11111111 11111111 Add 7F00-mask
// 00000000 00000000 00000000 00000011 movemask



// 3-Byte
// 00000000 00000010 00000010 00001010 Value
// 00000001 00000001 00000001 00000001 1's Mask
// 00000000 00000001 00000001 00000001 Min(value, 1's)
//
// 00000000 00000000 00000001 11111111 Pack(min)
// 00000001 00000001 00000001 00000001 1's Mask
// 00000000 00000000 00000001 00000001 16bit-min(packed-min)
//
// 01111111 00000000 01111111 00000000 7F00-mask
// 01111111 00000000 10000000 00000001 Add 7F00-mask
// 00000000 00000000 00000000 00000010 movemask


// 4-Byte
// 00000010 00000010 00000010 00001010 Value
// 00000001 00000001 00000001 00000001 1's Mask
// 00000001 00000001 00000001 00000001 Min(value, 1's)
//
// 00000000 00000000 11111111 11111111 Pack(min)
// 00000001 00000001 00000001 00000001 1's Mask
// 00000000 00000000 11111111 11111111 16bit-min(packed-min)
//
// 01111111 00000000 01111111 00000000 7F00-mask
// 01111111 00000000 11111111 11111111 Add 7F00-mask
// 00000000 00000000 00000000 00000011 movemask
// --


// ???
// 00000010 00000010 00000000 00001010 Value
// 00000001 00000001 00000001 00000001 1's Mask
// 00000001 00000001 00000000 00000001 Min(value, 1's)
//
// 00000000 00000000 11111111 00000001 Pack(min)
// 00000001 00000001 00000001 00000001 1's Mask
// 00000000 00000000 11111111 00000001 16bit-min(packed-min)
//
// 01111111 00000000 01111111 00000000 7F00-mask
// 01111111 00000000 11111111 11111111 Add 7F00-mask
// 00000000 00000000 00000000 00000011 movemask




// ???
// 10000000 10000000 10000000 10000000 Value
// 00000001 00000001 00000001 00000001 1's Mask
// 00000001 00000001 00000001 00000001 Min(value, 1's)
//
// 11111111 11111111 11111111 11111111 Pack(min)
// 00000001 00000001 00000001 00000001 1's Mask
// 00000000 00000000 11111111 00000001 16bit-min(packed-min)
//
// 01111111 00000000 01111111 00000000 7F00-mask
// 01111111 00000000 11111111 11111111 Add 7F00-mask
// 00000000 00000000 00000000 00000011 movemask