// Package stream provides logic that allows for you to operate on entire streams
// of uint32's or bytes.
package stream

const (
	HeaderMagic uint64 = 0x3fd76c17
	FooterMagic = ^HeaderMagic
)

type Opts uint64

const (
	// Delta indicates that the stream is delta encoded
	Delta = 1 << iota
	// Skip indicates that the stream contains a skip table that
	// can be used to skip around in the stream.
	Skip
)
