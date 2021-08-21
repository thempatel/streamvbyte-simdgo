package main

import (
	. "github.com/mmcloughlin/avo/build"
	"github.com/mmcloughlin/avo/operand"
	"github.com/mmcloughlin/avo/reg"
)

var (
	mask_01 uint16 = 0x1111
	mask_7F00 uint16 = 0x7F00
	lowerCtrl uint64 = 0xFF
	upperCtl = ^lowerCtrl

	name = "PutUint32x86_8"
	signature = "func(in []uint32, outBytes []byte, shuffle [][16]uint8, lenTable []uint8) (r uint16)"
)

func main() {
	TEXT(name, NOSPLIT, signature)
	Doc("PutUint32x86_8 encodes 8 32-bit unsigned integers at a time.")

	maskO1 := ConstData("mask_01", operand.U16(mask_01))
	mask7F00 := ConstData("mask_7F00", operand.U16(mask_7F00))

	onesMask := XMM()
	sevenFzerozero := XMM()
	VPBROADCASTW(maskO1, onesMask)
	VPBROADCASTW(mask7F00, sevenFzerozero)

	arrBase := operand.Mem{Base: Load(Param("in").Base(), GP64())}

	firstFour := XMM()
	secondFour := XMM()
	// load first 4 uint32's
	VLDDQU(arrBase, firstFour)
	// load second 4 uint32's
	VLDDQU(arrBase.Offset(16), secondFour)

	minFirstFour := XMM()
	minSecondFour := XMM()
	VPMINUB(onesMask, firstFour, minFirstFour)
	VPMINUB(onesMask, secondFour, minSecondFour)

	// Re-use minFirstFour register
	VPACKUSWB(minSecondFour, minFirstFour, minFirstFour)
	VPMINSW(onesMask, minFirstFour, minFirstFour)
	VPADDUSW(sevenFzerozero, minFirstFour, minFirstFour)

	ctrl := GP32()
	VPMOVMSKB(minFirstFour, ctrl)
	Store(ctrl.As16(), Return("r"))

	shuffleBase := Load(Param("shuffle").Base(), GP64())

	// Gives the index into the shuffle table for the first 4 numbers encoded
	// Move the lower 8 bytes into the register
	MOVBQZX(ctrl.As8(), reg.RAX)
	// Left shift by 4 to get the byte level offset for the shuffle table
	SHLQ(operand.Imm(4), reg.RAX)
	ADDQ(shuffleBase, reg.RAX)
	firstShuffle := XMM()
	VLDDQU(operand.Mem{Base: reg.RAX}, firstShuffle)

	MOVWQZX(ctrl.As16(), reg.RAX)
	SHRQ(operand.Imm(8), reg.RAX)
	SHLQ(operand.Imm(4), reg.RAX)
	ADDQ(shuffleBase, reg.RAX)
	secondShuffle := XMM()
	VLDDQU(operand.Mem{Base: reg.RAX}, secondShuffle)

	VPSHUFB(firstShuffle, firstFour, firstFour)
	VPSHUFB(secondShuffle, secondFour, secondFour)

	baseAddr := Load(Param("outBytes").Base(), GP64())
	VMOVDQU(firstShuffle, operand.Mem{Base: baseAddr})

	Load(Param("lenTable").Base(), reg.RAX)
	ANDL(operand.Imm(0xff), ctrl)
	ADDQ(ctrl.As64(), reg.RAX)

	length := GP8()
	MOVB(operand.Mem{Base: reg.RAX}, length)
	ADDQ(length.As64(), baseAddr)
	VMOVDQU(secondShuffle, operand.Mem{Base: baseAddr})

	RET()
	Generate()
}