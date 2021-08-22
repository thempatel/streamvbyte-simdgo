package encode

type Put8Impl func([]uint32, []byte) uint16
type cpuCheck func() bool

var (
	putImpl Put8Impl
)

func init() {
	if check() {
		putImpl = put8uint32
	} else {
		putImpl = put8uint32Scalar
	}
}

func Put8uint32(in []uint32, out []byte) uint16 {
	return putImpl(in, out)
}

func put8uint32Scalar(in []uint32, out []byte ) uint16 {

}