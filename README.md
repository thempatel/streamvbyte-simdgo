# Stream VByte SIMD Go

![Tests](https://github.com/theMPatel/streamvbyte-simdgo/actions/workflows/default.yaml/badge.svg)

This is a repository that contains a port of Stream VByte to Go. Notably, this repo takes extra care
to leverage SIMD techniques to achieve better performance. Currently, there is support for x86_64 architectures
that have AVX and AVX2 hardware instructions. In cases where that is not available, or on non x86_64 architectures
there is a portable scalar implementation. We also perform a runtime check to make sure that the necessary
ISA is available and if not fallback to the scalar approach.

There are several existing implementations:

1. [Reference C/C++](https://github.com/lemire/streamvbyte)
2. [Rust](https://bitbucket.org/marshallpierce/stream-vbyte-rust)
3. [Go](https://github.com/nelz9999/stream-vbyte-go)
   * Note: only has a scalar implementation which prompted this implementation with SIMD techniques.

## Benchmarks

```text
goos: darwin
goarch: amd64
pkg: github.com/theMPatel/streamvbyte-simdgo/pkg
cpu: Intel(R) Core(TM) i7-8700B CPU @ 3.20GHz
--
BenchmarkMemCopy8Uint32-12    	448033621	         2.613 ns/op	12247.79 MB/s

goos: darwin
goarch: amd64
pkg: github.com/theMPatel/streamvbyte-simdgo/pkg/decode
cpu: Intel(R) Core(TM) i7-8700B CPU @ 3.20GHz
--
BenchmarkGet8uint32Fast-12           	368871128	         3.234 ns/op	9895.19 MB/s
BenchmarkGet8uint32DeltaFast-12      	306617568	         3.853 ns/op	8304.19 MB/s
BenchmarkGet8uint32Scalar-12         	67530424	        18.59 ns/op	1721.51 MB/s
BenchmarkGet8uint32DeltaScalar-12    	65821027	        19.17 ns/op	1668.90 MB/s
BenchmarkGet8uint32Varint-12         	26296116	        46.45 ns/op	 688.88 MB/s
BenchmarkGet8uint32DeltaVarint-12    	21070124	        57.95 ns/op	 552.16 MB/s

goos: darwin
goarch: amd64
pkg: github.com/theMPatel/streamvbyte-simdgo/pkg/encode
cpu: Intel(R) Core(TM) i7-8700B CPU @ 3.20GHz
--
BenchmarkPut8uint32Fast-12           	310744244	         3.883 ns/op	8241.67 MB/s
BenchmarkPut8uint32DeltaFast-12      	276355396	         4.305 ns/op	7434.07 MB/s
BenchmarkPut8uint32Scalar-12         	41064336	        29.89 ns/op	1070.47 MB/s
BenchmarkPut8uint32DeltaScalar-12    	40731841	        29.17 ns/op	1096.87 MB/s
BenchmarkPut8uint32Varint-12         	48454465	        22.42 ns/op	1427.45 MB/s
BenchmarkPut8uint32DeltaVarint-12    	68148270	        22.23 ns/op	1439.20 MB/s

goos: darwin
goarch: amd64
pkg: github.com/theMPatel/streamvbyte-simdgo/pkg/stream/reader
cpu: Intel(R) Core(TM) i7-8700B CPU @ 3.20GHz
--
BenchmarkReadAllFast/Count_1e0-12         	100000000	        11.64 ns/op	 343.50 MB/s
BenchmarkReadAllFast/Count_1e1-12         	27834266	        43.15 ns/op	 927.09 MB/s
BenchmarkReadAllFast/Count_1e2-12         	11309079	       106.0 ns/op	3773.54 MB/s
BenchmarkReadAllFast/Count_1e3-12         	 1673222	       720.6 ns/op	5550.61 MB/s
BenchmarkReadAllFast/Count_1e4-12         	  171164	      6941 ns/op	5763.09 MB/s
BenchmarkReadAllFast/Count_1e5-12         	   17080	     70616 ns/op	5664.44 MB/s
BenchmarkReadAllFast/Count_1e6-12         	    1640	    709504 ns/op	5637.74 MB/s
BenchmarkReadAllFast/Count_1e7-12         	     154	   7670283 ns/op	5214.93 MB/s
BenchmarkFastRead-12                      	  421057	      2885 ns/op	5678.46 MB/s
BenchmarkReadAllScalar/Count_1e0-12       	128565745	         9.355 ns/op	 427.59 MB/s
BenchmarkReadAllScalar/Count_1e1-12       	37689912	        32.12 ns/op	1245.18 MB/s
BenchmarkReadAllScalar/Count_1e2-12       	 4740482	       252.9 ns/op	1581.38 MB/s
BenchmarkReadAllScalar/Count_1e3-12       	  482290	      2492 ns/op	1605.11 MB/s
BenchmarkReadAllScalar/Count_1e4-12       	   17554	     68111 ns/op	 587.28 MB/s
BenchmarkReadAllScalar/Count_1e5-12       	    1534	    776822 ns/op	 514.92 MB/s
BenchmarkReadAllScalar/Count_1e6-12       	     153	   7792223 ns/op	 513.33 MB/s
BenchmarkReadAllScalar/Count_1e7-12       	      14	  78339118 ns/op	 510.60 MB/s

goos: darwin
goarch: amd64
pkg: github.com/theMPatel/streamvbyte-simdgo/pkg/stream/writer
cpu: Intel(R) Core(TM) i7-8700B CPU @ 3.20GHz
--
BenchmarkWriteAllFast/Count_1e0-12         	54667987	        21.98 ns/op	 181.97 MB/s
BenchmarkWriteAllFast/Count_1e1-12         	28875541	        41.76 ns/op	 957.95 MB/s
BenchmarkWriteAllFast/Count_1e2-12         	 7327003	       161.7 ns/op	2473.45 MB/s
BenchmarkWriteAllFast/Count_1e3-12         	  992776	      1229 ns/op	3255.02 MB/s
BenchmarkWriteAllFast/Count_1e4-12         	   98493	     11569 ns/op	3457.37 MB/s
BenchmarkWriteAllFast/Count_1e5-12         	   10000	    108539 ns/op	3685.30 MB/s
BenchmarkWriteAllFast/Count_1e6-12         	     915	   1197473 ns/op	3340.37 MB/s
BenchmarkWriteAllFast/Count_1e7-12         	     106	  10155791 ns/op	3938.64 MB/s
BenchmarkFastWrite-12                      	  247346	      4677 ns/op	3503.21 MB/s
BenchmarkWriteAllScalar/Count_1e0-12       	57365708	        21.07 ns/op	 189.83 MB/s
BenchmarkWriteAllScalar/Count_1e1-12       	18698986	        63.59 ns/op	 629.05 MB/s
BenchmarkWriteAllScalar/Count_1e2-12       	 2709844	       440.5 ns/op	 908.03 MB/s
BenchmarkWriteAllScalar/Count_1e3-12       	  286054	      4162 ns/op	 960.98 MB/s
BenchmarkWriteAllScalar/Count_1e4-12       	   12038	     99717 ns/op	 401.14 MB/s
BenchmarkWriteAllScalar/Count_1e5-12       	    1141	   1049287 ns/op	 381.21 MB/s
BenchmarkWriteAllScalar/Count_1e6-12       	     100	  10509348 ns/op	 380.61 MB/s
BenchmarkWriteAllScalar/Count_1e7-12       	      10	 105091277 ns/op	 380.62 MB/s
```

A note on the benchmarks: An array of random uint32's is generated and then encoded/decoded over
and over again. An attempt is made to ensure that some of these benchmarks reflect the most probable
real world performance metrics.

---
Stream VByte uses the same underlying format as Google's Group Varint approach. Lemire et al. wanted
to see if there was a way to improve the performance even more and introduced a clever twist to enable
better performance via SIMD techniques. The basic goal of the Group Varint format is to be able to
achieve similar compression characteristics as the VByte format for integers and also be able to load
and process them really quickly.

## VByte format

The insight that backs the VByte encoding is noticing that you oftentimes don't need 32 bits to
encode a 32-bit integer. Take for example an unsigned integer that is less than 2^8 (256). This
integer will have bits set in the lowest byte of a 32-bit integer, while the remaining 3 bytes will
simply be zeros.

```
111 in binary:

00000000 00000000 00000000 01101111
```

An approach you can take to compress this integer is to encode the integer using a variable
number of bytes. For example, you can use the lower 7 bits to store data, i.e. bits
from the original integer, and then use the MSB as a continuation bit. If the MSB bit is on, i.e.
is 1, then more bytes are needed to decode this particular integer. Below is an example where
you might need 2 bytes to store the number 1234.

```
1234 in binary:

00000000 00000000 00000100 11010010

Num compressed:

v          v          Continuation bits
0|0001001| 1|1010010|
    ^           ^     Data bits
```

If you want to decode this integer, you simply build up the number iteratively. I.e. you OR the
last 7 bits of every byte shifted to the appropriate length to your 32-bit integer until you
find a byte that doesn't have a continuation bit set. Note that this works the same for 64-bit
numbers.

The problem with this approach is that it can introduce a lot of branch mis-predictions during encoding/decoding.
During the decoding phase, you don't know ahead of time the number of bytes that were used to encode the integer
you are currently processing and so you need to iterate until you find a byte without a continuation bit on.
If you have integers that are nonuniform, i.e. integers that require random numbers of bytes to encode relative
to one another, this can pose a challenge to the processor's branch predictor. These mis-predictions can cause
major slowdowns in processor pipelines and so was born the Group Varint format.

## Group Varint format

The Group Varint (varint-GB) format assumes that everything you hope to achieve, you can do with 32-bit integers.
It introduces the concept of a control byte which is simply a byte that stores the encoded
lengths of a group of 4 32-bit integers, hence Group Varint. 32-bit integers only require up to 4 bytes
to properly encode. This means that you can represent their lengths with 2 bits using a zero-indexed length
i.e. 0, 1, 2, and 3 to represent integers that require 1, 2, 3 and 4 bytes to encode, respectively.

```
00000000 00000000 00000000 01101111  =        111 
00000000 00000000 00000100 11010010  =       1234
00000000 00001100 00001010 10000011  =     789123
01000000 00000000 00000000 00000000  = 1073741824

Num         Len      2-bit control
----------------------------------
111          1                0b00 
1234         2                0b01
789123       3                0b10
1073741824   4                0b11

Final Control byte
0b11100100

Encoded data (little endian right-to-left bottom-to-top) 
0b01000000 0b00000000 0b00000000 0b00000000 0b00001100
0b00001010 0b10000011 0b00000100 0b11010010 0b01101111
```

You can then prefix every group of 4 encoded 32-bit integers with their control byte and then use it during decoding.
The obvious downside is that you pay a storage cost of one byte for every 4 integers you want to encode. For 2^20 
encoded integers, that's an extra 256 KB of extra space: totally marginal. The great upside, though, is that
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
	//  	_ = b[3]
	//  	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
	// }
	return binary.LittleEndian.Uint32(buf)
}

func main() {
	ctrl := uint8(0b11_10_01_00)
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
	
	_ = decodeOne(data, len0)                   // 111
	_ = decodeOne(data[len0:], len1)            // 1234
	_ = decodeOne(data[len0+len1:], len2)       // 789_123
	_ = decodeOne(data[len0+len1+len2:], len3)  // 1_073_741_824
}
```

## Stream VByte format

Unfortunately, accelerating decoding of the varint-GB format with only SIMD techniques
has proven unsuccessful. The below excerpt from the paper outlines why. 

> To understand why it might be difficult to accelerate the decoding of data compressed in the VARINT-GB
> format compared to the VARINT-G8IU format, consider that we cannot decode faster than we can access the
> control bytes. In VARINT-G8IU, the control bytes are conveniently always located nine compressed bytes
> apart. Thus while a control byte is being processed, or even before, our superscalar processor can load
> and start processing upcoming control bytes, as their locations are predictable. Instructions depending
> on these control bytes can be reordered by the processor for best performance. However, in the VARINT-GB
> format, there is a strong data dependency: the location of the next control byte depends on the current
> control byte. This increases the risk that the processor remains underutilized, delayed by the latency
> between issuing the load for the next control byte and waiting for it to be ready.

Additionally, they prove that decoding 4 integers at a time using 128-bit registers is faster than trying
to decode a variable number of integers that fit into an 8-byte register, i.e. the varint-G8IU approach.

### SIMD control byte generation algorithm

Lemire et al. have devised a brilliant SIMD algorithm for simultaneously generating two control bytes
for a group of 8 integers. The best way to understand this algorithm is to understand how it works on 
a single integer and then assume it works in a vectorized form (it does). Going forward we'll use
*control bits stream* to represent these control bytes we are building. 

```
00000000 00000000 00000100 11010010 // 1234
```

Let's take one of the previous integers that we were looking at, `1234`, and walk through an example
of how the 2-bit control is generated for it using SIMD techniques. The goal is to be able to, for
any 32-bit integer, generate a 2-bit zero indexed length value. For example, if you have an integer
that requires 2 bytes to be encoded, we want for the algorithm to generate `0b01`.

```
00000000 00000000 00000100 11010010 // 1234
00000001 00000001 00000001 00000001 // 0x0101 mask
----------------------------------- // byte-min(1234, 0x0101)
00000000 00000000 00000001 00000001
```

The algorithm first uses a mask where every byte is equal to 1. If you perform a per-byte min operation
on our integer and the 1's mask, the result will have a 1 at every byte that had a value in the original
integer. 

```
00000000 00000000 00000001 00000001
----------------------------------- // pack unsigned saturating 16-bit to 8-bit
00000000 00000000 00000000 11111111
```

Now you perform a 16-bit to 8-bit unsigned saturating pack operation. Practically this means that you're
taking every 16-bit value and trying to shove that into 8 bits. If the 16-bit integer is larger than
the largest unsigned integer 8 bits can support, the pack saturates to the largest unsigned 8-bit value. 

Why this is performed will become more clear in the subsequent steps, however, at a high level, for every
integer you want to encode, you want for the MSB of two consecutive bytes in the control bits stream
to be representative of the final 2-bit control. For example, if you have a 3-byte integer, you want the
MSB of two consecutive bytes to be 1 and 0, in that order. The reason you would want this is that
there is a vector pack instruction that takes the MSB from every byte in the control bits stream
and packs it into the lowest byte. This would thus represent the value `0b10` in the final byte for
this 3-byte integer, which is what we want.

Performing a 16-bit to 8-bit unsigned saturating pack has the effect that you can use the saturation
behavior to conditionally turn on the MSB of these bytes depending on which bytes have values in the
original 32-bit integer.

```
00000000 00000000 00000000 11111111 // control bits stream
00000001 00000001 00000001 00000001 // 0x0101 mask
----------------------------------- // signed 16-bit min
00000000 00000000 00000000 11111111
```

We then take the 1's mask we used before and perform a __signed 16-bit__ min operation. The reason for this
is more clear if you look at an example using a 3-byte integer.

```
00000000 00001100 00001010 10000011 // 789123
00000001 00000001 00000001 00000001 // 0x0101 mask
----------------------------------- // byte-min(789123, 0x0101)
00000000 00000001 00000001 00000001
----------------------------------- // pack unsigned saturating 16-bit to 8-bit
00000000 00000000 00000001 11111111
00000001 00000001 00000001 00000001 // 0x0101 mask
----------------------------------- // signed 16-bit min
00000000 00000000 00000001 00000001
```

The signed 16-bit min operation has three important effects.

First, for 3-byte integers, it has the effect of turning off the MSB of the lowest byte. This is necessary
because a 3-byte integer should have a 2-bit control that is `0b10` and without this step using the MSB pack
operation would result in a 2-bit control that looks something like `0b_1`, where the lowest bit is on.
Obviously this is wrong, since only integers that require 2 or 4 bytes to encode should have that lower bit
on, i.e. 1 or 3 as a zero-indexed length.

Second, for 4-byte integers, the signed aspect has the effect of leaving both MSBs of the 2 bytes on. When using the
MSB pack operation later on, it will result in a 2-bit control value of `0b11`, which is what we want.

Third, for 1 and 2 byte integers, it has no effect. This is great for 2-byte values since the MSB will remain on
and 1 byte values will not have any MSB on anyways, so it is effectively a noop in both scenarios.

```
00000000 00000000 00000000 11111111 // control bits stream (original 1234)
01111111 00000000 01111111 00000000 // 0x7F00 mask
----------------------------------- // add unsigned saturating 16-bit
01111111 00000000 01111111 11111111
```

Next, we take a mask with the value `0x7F00` and perform an unsigned saturating add to the control bits stream.
In the case for the integer `1234` this has no real effect. We maintain the MSB in the lowest byte. You'll note,
however, that the only byte that has its MSB on is the last one, so performing an MSB pack operation would result
in a value of `0b0001`, which is what we want. An example of this step on the integer `789123` might paint a clearer
picture.

```
00000000 00000000 00000001 00000001 // control bits stream (789123)
01111111 00000000 01111111 00000000 // 0x7F00 mask
----------------------------------- // add unsigned saturating 16-bit
01111111 00000000 11111111 00000001
```

You'll note here that the addition of `0x01` with `0x7F` in the upper byte results in the MSB of the resulting upper
byte turning on. The MSB in the lower byte remains off and now an MSB pack operation will resolve to `0b0010`,
which is what we want. The unsigned saturation behavior is really important for 4-byte numbers that only have
bits in the most significant byte on. An example below:

```
01000000 00000000 00000000 00000000 // 1073741824
00000001 00000001 00000001 00000001 // 0x0101 mask
----------------------------------- // byte-min(1073741824, 0x0101)
00000001 00000000 00000000 00000000
----------------------------------- // pack unsigned saturating 16-bit to 8-bit
00000000 00000000 11111111 00000000
00000001 00000001 00000001 00000001 // 0x0101 mask
----------------------------------- // signed 16-bit min
00000000 00000000 11111111 00000000
01111111 00000000 01111111 00000000 // 0x7F00 mask
----------------------------------- // add unsigned saturating 16-bit
01111111 00000000 11111111 11111111
```

Note here that because only the upper byte had a value in it, the lowest byte in the control bits stream remains
zero for the duration of the algorithm. This poses an issue, since for a 4-byte value, we want for the 2-bit
control to result in a value of `0b11`. Performing a 16-bit unsigned *saturating* addition has the effect of
turning on all bits in the lower byte, and thus we get a result with the MSB in the lower byte on. 

```
01111111 00000000 11111111 00000001 // control bits stream (789123)
----------------------------------- // move byte mask 
00000000 00000000 00000000 00000010 // 2-bit control 
```

The final move byte mask is performed on the control bits stream, and we now have the result we wanted. Now that you
see that this works for 1 integer, you know how it can work for 8 integers simultaneously, since we use vector
instructions that operate on 128 bit registers.

### SIMD integer packing/unpacking

The next problem to be solved is how to take a group of 4 integers, and compress it by removing extraneous/unused
bytes so that all you're left with is a stream of data bytes with real information. Let's take two numbers from
our examples above.

```
               789123                                 1234
00000000 00001100 00001010 10000011 | 00000000 00000000 00000100 11010010
-------------------------------------------------------------------------
         00001100 00001010 10000011   00000100 11010010      // packed
```

Here, we can use a shuffle operation. Vector shuffle operations rearrange the bytes in an input register according
to some provided mask into a destination register. Every position in the mask stores an offset into the source
vector stream that represents the data byte that should go into that position.

```
input [1234, 789123] (little endian R-to-L)
00000000 00001100 00001010 10000011 00000000 00000000 00000100 11010010
            |       |         |                             |        |
            |       |         |____________________         |        |
            |       |_____________________         |        |        |
            |____________________         |        |        |        |
                                 v        v        v        v        v
    0xff     0xff     0xff     0x06     0x05     0x04     0x01     0x00 // mask in hex
-----------------------------------------------------------------------
00000000 00000000 00000000 00001100 00001010 10000011 00000100 11010010 // packed
```

We keep a prebuilt lookup table that contains a mapping from control byte to the necessary mask and simply
load that after we construct the control byte above. In addition, we keep a lookup table for a mapping from
control bytes to total encoded length. This allows us to know by how much to increment the output pointer and
overwrite, for example, the redundant upper 3 bytes in the above shuffle example.

Unpacking during decoding is the same as the above, but in reverse. We need to go from a packed format
to an unpacked memory format. We keep lookup tables to maintain a mapping from control byte to the reverse
shuffle mask, and then perform a shuffle operation to output to an `uint32` array.

# References

[Stream VByte: Faster Byte-Oriented Integer Compression](https://arxiv.org/pdf/1709.08990.pdf)
