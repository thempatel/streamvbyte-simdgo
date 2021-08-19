package main

import (
	. "github.com/mmcloughlin/avo/build"
	"github.com/mmcloughlin/avo/operand"
)

var (
	mask_01 uint16 = 0x1111
	mask_7F00 uint16 = 0x7F00
)

func main() {
	TEXT("x86ControlBytes8", NOSPLIT, "func(in []uint32) (r uint32)")
	Doc("Encodes 8 32-bit unsigned integers at a time.")

	a := ConstData("mask_01", operand.U16(mask_01))
	b := ConstData("mask_7F00", operand.U16(mask_7F00))
	onesMask := XMM()
	sevenFzerozero := XMM()

	VPBROADCASTW(a, onesMask)
	VPBROADCASTW(b, sevenFzerozero)

	firstHalf := operand.Mem{Base: Load(Param("in").Base(), GP64())}

	r0 := XMM()
	r1 := XMM()

	VLDDQU(firstHalf, r0) // load first 4
	VLDDQU(firstHalf.Offset(16), r1)

	PMINUB(onesMask, r0)
	PMINUB(onesMask, r1)

	PACKUSWB(r1, r0)
	PMINSW(onesMask, r0)
	PADDUSW(sevenFzerozero, r0)

	dest := GP32()
	PMOVMSKB(r0, dest)
	Store(dest, Return("r"))
	RET()
	Generate()
}