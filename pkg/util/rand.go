package util

// This file provides a more uniform random number generator that creates
// numbers that have a more normal distribution along the number of bytes
// required to encode them. This is needed because the larger encoded bytes
// i.e. 3 and 4 bytes have more numbers to pick from versus those that require
// just 1 or 2. Thus, using a normally generated number is more likely to produce
// a number that requires 3 or 4 bytes to encode.

import (
	"math"
	"math/rand"
)

type generator func() uint32

func randUint32Range(low, high uint32) uint32 {
	return (rand.Uint32() % (high - low + 1)) + low
}

var (
	generators = []generator{
		// 1 byte,
		func() uint32 {
			return randUint32Range(0, 1<<8)
		},
		// 2 byte,
		func() uint32 {
			return randUint32Range(1<<8, 1<<16)
		},
		// 3 byte,
		func() uint32 {
			return randUint32Range(1<<16, 1<<24)
		},
		// 4 byte,
		func() uint32 {
			return randUint32Range(1<<24, math.MaxUint32)
		},
	}
)

// RandUint32 generates a random number that is also uniformly random
// on the axis for the number of bytes required to encode it. It first
// randomly chooses a byte length, i.e. 1, 2, 3 or 4 and then randomly
// generates a number whose encoded length would be that length.
func RandUint32() uint32 {
	return generators[rand.Int()%4]()
}
