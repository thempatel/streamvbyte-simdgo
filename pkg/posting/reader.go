package posting

type Reader interface {
	// Read will read a count number of uint32's from the input
	// stream constructed at the outset
	Read(count int) []uint32
}
