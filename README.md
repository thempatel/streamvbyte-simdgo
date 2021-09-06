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
MemCopy8Uint32-12    	463986302	         2.608 ns/op	12269.03 MB/s

goos: darwin
goarch: amd64
pkg: github.com/theMPatel/streamvbyte-simdgo/pkg/decode
cpu: Intel(R) Core(TM) i7-8700B CPU @ 3.20GHz
--
Get8uint32Fast-12           	377839186	         3.170 ns/op	10095.99 MB/s
Get8uint32DeltaFast-12      	298522095	         4.455 ns/op	7183.20 MB/s
Get8uint32Scalar-12         	63384603	        19.28 ns/op	1659.36 MB/s
Get8uint32DeltaScalar-12    	58705828	        20.04 ns/op	1596.46 MB/s
Get8uint32Varint-12         	27369775	        43.77 ns/op	 731.10 MB/s
Get8uint32DeltaVarint-12    	20924770	        57.30 ns/op	 558.46 MB/s

goos: darwin
goarch: amd64
pkg: github.com/theMPatel/streamvbyte-simdgo/pkg/encode
cpu: Intel(R) Core(TM) i7-8700B CPU @ 3.20GHz
--
Put8uint32Fast-12           	297620898	         3.864 ns/op	8281.18 MB/s
Put8uint32DeltaFast-12      	276545827	         4.350 ns/op	7356.59 MB/s
Put8uint32Scalar-12         	41200776	        28.59 ns/op	1119.30 MB/s
Put8uint32DeltaScalar-12    	37773458	        30.65 ns/op	1044.11 MB/s
Put8uint32Varint-12         	58668867	        17.20 ns/op	1860.67 MB/s
Put8uint32DeltaVarint-12    	61446153	        22.88 ns/op	1398.80 MB/s

goos: darwin
goarch: amd64
pkg: github.com/theMPatel/streamvbyte-simdgo/pkg/stream/reader
cpu: Intel(R) Core(TM) i7-8700B CPU @ 3.20GHz
--
ReadAllFast/Count_1e0-12 	99354789	        12.24 ns/op	 326.80 MB/s
ReadAllFast/Count_1e1-12 	28076071	        42.81 ns/op	 934.43 MB/s
ReadAllFast/Count_1e2-12 	11041639	       107.2 ns/op	3730.16 MB/s
ReadAllFast/Count_1e3-12 	 1645387	       729.9 ns/op	5480.00 MB/s
ReadAllFast/Count_1e4-12 	  170894	      7034 ns/op	5686.52 MB/s
ReadAllFast/Count_1e5-12 	   16848	     70969 ns/op	5636.29 MB/s
ReadAllFast/Count_1e6-12 	    1513	    728516 ns/op	5490.62 MB/s
ReadAllFast/Count_1e7-12 	     152	   7835111 ns/op	5105.22 MB/s
ReadAllDeltaFast/Count_1e0-12         	92727970	        13.10 ns/op	 305.44 MB/s
ReadAllDeltaFast/Count_1e1-12         	26164140	        45.89 ns/op	 871.61 MB/s
ReadAllDeltaFast/Count_1e2-12         	 9458992	       128.5 ns/op	3113.55 MB/s
ReadAllDeltaFast/Count_1e3-12         	 1277408	       934.4 ns/op	4280.69 MB/s
ReadAllDeltaFast/Count_1e4-12         	  144405	      8318 ns/op	4808.88 MB/s
ReadAllDeltaFast/Count_1e5-12         	   14444	     83151 ns/op	4810.55 MB/s
ReadAllDeltaFast/Count_1e6-12         	    1426	    846305 ns/op	4726.43 MB/s
ReadAllDeltaFast/Count_1e7-12         	     127	   9337355 ns/op	4283.87 MB/s
ReadAllScalar/Count_1e0-12            	122650209	         9.770 ns/op	 409.43 MB/s
ReadAllScalar/Count_1e1-12            	38012136	        31.63 ns/op	1264.64 MB/s
ReadAllScalar/Count_1e2-12            	 4999376	       241.6 ns/op	1655.30 MB/s
ReadAllScalar/Count_1e3-12            	  500337	      2459 ns/op	1626.38 MB/s
ReadAllScalar/Count_1e4-12            	   50247	     24034 ns/op	1664.34 MB/s
ReadAllScalar/Count_1e5-12            	    5032	    238354 ns/op	1678.17 MB/s
ReadAllScalar/Count_1e6-12            	     499	   2405669 ns/op	1662.74 MB/s
ReadAllScalar/Count_1e7-12            	      46	  24533207 ns/op	1630.44 MB/s
ReadAllDeltaScalar/Count_1e0-12       	100000000	        10.32 ns/op	 387.49 MB/s
ReadAllDeltaScalar/Count_1e1-12       	36915704	        32.52 ns/op	1230.08 MB/s
ReadAllDeltaScalar/Count_1e2-12       	 4818140	       249.8 ns/op	1601.58 MB/s
ReadAllDeltaScalar/Count_1e3-12       	  512492	      2374 ns/op	1685.20 MB/s
ReadAllDeltaScalar/Count_1e4-12       	   51004	     23639 ns/op	1692.11 MB/s
ReadAllDeltaScalar/Count_1e5-12       	    3568	    333168 ns/op	1200.60 MB/s
ReadAllDeltaScalar/Count_1e6-12       	     520	   2304864 ns/op	1735.46 MB/s
ReadAllDeltaScalar/Count_1e7-12       	      48	  24810555 ns/op	1612.22 MB/s
ReadAllVarint/Count_1e0-12            	121348074	         9.967 ns/op	 401.34 MB/s
ReadAllVarint/Count_1e1-12            	21056739	        57.34 ns/op	 697.64 MB/s
ReadAllVarint/Count_1e2-12            	 2025081	       589.0 ns/op	 679.15 MB/s
ReadAllVarint/Count_1e3-12            	  205881	      5851 ns/op	 683.69 MB/s
ReadAllVarint/Count_1e4-12            	   20906	     57446 ns/op	 696.31 MB/s
ReadAllVarint/Count_1e5-12            	    2037	    580620 ns/op	 688.92 MB/s
ReadAllVarint/Count_1e6-12            	     208	   5755083 ns/op	 695.04 MB/s
ReadAllVarint/Count_1e7-12            	      20	  57872736 ns/op	 691.17 MB/s
ReadAllDeltaVarint/Count_1e0-12       	139763250	         8.318 ns/op	 480.87 MB/s
ReadAllDeltaVarint/Count_1e1-12       	19199100	        62.49 ns/op	 640.11 MB/s
ReadAllDeltaVarint/Count_1e2-12       	 2149660	       556.6 ns/op	 718.65 MB/s
ReadAllDeltaVarint/Count_1e3-12       	  207122	      5810 ns/op	 688.41 MB/s
ReadAllDeltaVarint/Count_1e4-12       	   22680	     53200 ns/op	 751.88 MB/s
ReadAllDeltaVarint/Count_1e5-12       	    2145	    500177 ns/op	 799.72 MB/s
ReadAllDeltaVarint/Count_1e6-12       	     228	   5262741 ns/op	 760.06 MB/s
ReadAllDeltaVarint/Count_1e7-12       	      27	  42000722 ns/op	 952.36 MB/s

goos: darwin
goarch: amd64
pkg: github.com/theMPatel/streamvbyte-simdgo/pkg/stream/writer
cpu: Intel(R) Core(TM) i7-8700B CPU @ 3.20GHz
--
WriteAllFast/Count_1e0-12 	54152408	        22.05 ns/op	 181.36 MB/s
WriteAllFast/Count_1e1-12 	27681948	        43.17 ns/op	 926.49 MB/s
WriteAllFast/Count_1e2-12 	 7136480	       167.0 ns/op	2395.79 MB/s
WriteAllFast/Count_1e3-12 	  928952	      1273 ns/op	3141.14 MB/s
WriteAllFast/Count_1e4-12 	   96117	     12012 ns/op	3329.93 MB/s
WriteAllFast/Count_1e5-12 	    9718	    114260 ns/op	3500.80 MB/s
WriteAllFast/Count_1e6-12 	     879	   1242927 ns/op	3218.21 MB/s
WriteAllFast/Count_1e7-12 	     100	  10368754 ns/op	3857.74 MB/s
WriteAllDeltaFast/Count_1e0-12         	50489378	        23.38 ns/op	 171.06 MB/s
WriteAllDeltaFast/Count_1e1-12         	26866423	        45.03 ns/op	 888.33 MB/s
WriteAllDeltaFast/Count_1e2-12         	 6695125	       175.8 ns/op	2275.37 MB/s
WriteAllDeltaFast/Count_1e3-12         	  899895	      1391 ns/op	2875.71 MB/s
WriteAllDeltaFast/Count_1e4-12         	   90394	     12958 ns/op	3086.82 MB/s
WriteAllDeltaFast/Count_1e5-12         	   10000	    122319 ns/op	3270.13 MB/s
WriteAllDeltaFast/Count_1e6-12         	     945	   1249546 ns/op	3201.16 MB/s
WriteAllDeltaFast/Count_1e7-12         	     100	  11461852 ns/op	3489.84 MB/s
WriteAllScalar/Count_1e0-12            	56106489	        21.72 ns/op	 184.18 MB/s
WriteAllScalar/Count_1e1-12            	18309972	        65.09 ns/op	 614.51 MB/s
WriteAllScalar/Count_1e2-12            	 2776918	       433.5 ns/op	 922.63 MB/s
WriteAllScalar/Count_1e3-12            	  289309	      4209 ns/op	 950.38 MB/s
WriteAllScalar/Count_1e4-12            	   29497	     40884 ns/op	 978.38 MB/s
WriteAllScalar/Count_1e5-12            	    3027	    399959 ns/op	1000.10 MB/s
WriteAllScalar/Count_1e6-12            	     296	   4010161 ns/op	 997.47 MB/s
WriteAllScalar/Count_1e7-12            	      28	  38753790 ns/op	1032.16 MB/s
WriteAllDeltaScalar/Count_1e0-12       	54981757	        21.90 ns/op	 182.65 MB/s
WriteAllDeltaScalar/Count_1e1-12       	17823349	        67.10 ns/op	 596.14 MB/s
WriteAllDeltaScalar/Count_1e2-12       	 2711672	       442.4 ns/op	 904.09 MB/s
WriteAllDeltaScalar/Count_1e3-12       	  292664	      4130 ns/op	 968.62 MB/s
WriteAllDeltaScalar/Count_1e4-12       	   29340	     41014 ns/op	 975.28 MB/s
WriteAllDeltaScalar/Count_1e5-12       	    2289	    516113 ns/op	 775.02 MB/s
WriteAllDeltaScalar/Count_1e6-12       	     302	   3930860 ns/op	1017.59 MB/s
WriteAllDeltaScalar/Count_1e7-12       	      30	  41357670 ns/op	 967.17 MB/s
WriteAllVarint/Count_1e0-12            	208214545	         5.720 ns/op	 699.32 MB/s
WriteAllVarint/Count_1e1-12            	43083270	        28.02 ns/op	1427.34 MB/s
WriteAllVarint/Count_1e2-12            	 4972045	       242.8 ns/op	1647.67 MB/s
WriteAllVarint/Count_1e3-12            	  499011	      2409 ns/op	1660.60 MB/s
WriteAllVarint/Count_1e4-12            	   51022	     23590 ns/op	1695.67 MB/s
WriteAllVarint/Count_1e5-12            	    5216	    231741 ns/op	1726.07 MB/s
WriteAllVarint/Count_1e6-12            	     518	   2305364 ns/op	1735.08 MB/s
WriteAllVarint/Count_1e7-12            	      50	  24905825 ns/op	1606.05 MB/s
WriteAllDeltaVarint/Count_1e0-12       	175269966	         6.792 ns/op	 588.93 MB/s
WriteAllDeltaVarint/Count_1e1-12       	51799438	        23.38 ns/op	1710.63 MB/s
WriteAllDeltaVarint/Count_1e2-12       	 5417458	       221.3 ns/op	1807.60 MB/s
WriteAllDeltaVarint/Count_1e3-12       	  539414	      2243 ns/op	1783.48 MB/s
WriteAllDeltaVarint/Count_1e4-12       	   52717	     22753 ns/op	1757.99 MB/s
WriteAllDeltaVarint/Count_1e5-12       	    5716	    210456 ns/op	1900.63 MB/s
WriteAllDeltaVarint/Count_1e6-12       	     495	   2453672 ns/op	1630.21 MB/s
WriteAllDeltaVarint/Count_1e7-12       	      70	  17491186 ns/op	2286.87 MB/s
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
