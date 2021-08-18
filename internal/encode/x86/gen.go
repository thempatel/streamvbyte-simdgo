package x86

//go:generate go run ./main/x86_encode.go -out ./x86_encode.s -stubs stub.go

func ControlBytes(in []uint32) uint32 {
	return x86ControlBytes8(in)
}
