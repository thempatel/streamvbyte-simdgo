package main

import (
	. "github.com/mmcloughlin/avo/build"
	"github.com/mmcloughlin/avo/operand"
	"github.com/mmcloughlin/avo/reg"
)

var (
	mask_01 uint8 = 0x11
	mask_7F00 uint16 = 0x7F00
)

func main() {
	TEXT("x86ControlBytes8", NOSPLIT, "func(in []uint32) uint32")
	Doc("Encodes 8 32-bit unsigned integers at a time.")

	a := GP8()
	b := GP16()
	MOVB(operand.U8(mask_01), a)
	MOVW(operand.U16(mask_7F00), b)

	onesMask := XMM()
	sevenFzerozero := XMM()

	VPBROADCASTB(operand.Mem{Base: a}, onesMask)
	VPBROADCASTW(operand.Mem{Base: b}, sevenFzerozero)

	firstHalf := operand.Mem{Base: Load(Param("in").Base(), GP64())}

	r0 := XMM()
	r1 := XMM()

	VLDDQU(firstHalf, r0) // load first 4
	VLDDQU(firstHalf.Offset(4*32), r1)

	PMINUB(onesMask, r0)
	PMINUB(onesMask, r1)

	PACKUSWB(r1, r0)
	PMINSW(onesMask, r0)
	PADDUSW(sevenFzerozero, r0)

	ANDL(operand.U32(0), reg.EAX)
	PMOVMSKB(r0, reg.EAX)
	Store(reg.EAX, ReturnIndex(0))
	RET()
	Generate()
}