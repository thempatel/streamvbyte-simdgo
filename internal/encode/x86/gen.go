package x86

//go:generate go run ./main/asm.go -out ./x86_encode.s -stubs stub.go
//go:generate go run ./main/asm2.go -out ./x86_getshuffle.s -stubs get_shuffle.go
