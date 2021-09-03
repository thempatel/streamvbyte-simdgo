package stream

type Reader interface {
	// ReadAll will read all uint32s from the input stream.
	ReadAll() []uint32
}
