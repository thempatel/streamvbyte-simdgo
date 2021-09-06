// +build !amd64

package reader

func ReadAllFast(count int, stream []byte, out []uint32) {
	panic("unreachable")
}
