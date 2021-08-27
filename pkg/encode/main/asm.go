package main

import (
	"fmt"

	. "github.com/mmcloughlin/avo/build"
	"github.com/mmcloughlin/avo/operand"
	"github.com/mmcloughlin/avo/reg"
)

const (
	mask1111 uint16 = 0x1111
	mask7F00 uint16 = 0x7F00

	mask1111Str = "mask1111"
	mask7F00Str = "mask7F00"

	name = "put8uint32Fast"
	pIn = "in"
	pOut = "outBytes"
	pShuffle = "shuffle"
	pLenTable = "lenTable"
	pR = "r"
)

var (
	signature = fmt.Sprintf(
		"func(%s []uint32, %s []byte, %s *[256][16]uint8, %s *[256]uint8) (%s uint16)",
		pIn, pOut, pShuffle, pLenTable, pR)
)

func main() {
	TEXT(name, NOSPLIT, signature)
	Doc("put8uint32Fast encodes 8 32-bit unsigned integers at a time.")

	mask1111R := ConstData(mask1111Str, operand.U16(mask1111))
	mask7F00R := ConstData(mask7F00Str, operand.U16(mask7F00))

	onesMask := XMM()
	sevenFzerozero := XMM()
	VPBROADCASTW(mask1111R, onesMask)
	VPBROADCASTW(mask7F00R, sevenFzerozero)

	arrBase := operand.Mem{Base: Load(Param(pIn).Base(), GP64())}

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
	Store(ctrl.As16(), Return(pR))

	shuffleBase := Load(Param(pShuffle), GP64())
	firstShuffle := loadCtrl16Shuffle(shuffleBase, ctrl, false)
	secondShuffle := loadCtrl16Shuffle(shuffleBase, ctrl, true)

	// TODO(milan): change to use memory operands
	VPSHUFB(firstShuffle, firstFour, firstFour)
	VPSHUFB(secondShuffle, secondFour, secondFour)

	firstAddr := Load(Param(pOut).Base(), GP64())
	secondAddr := GP64()
	MOVQ(firstAddr, secondAddr)

	lenValue := loadLenValue(ctrl)

	MOVBQZX(operand.Mem{Base: lenValue}, lenValue)
	ADDQ(lenValue, secondAddr)

	VMOVDQU(firstFour, operand.Mem{Base: firstAddr})
	VMOVDQU(secondFour, operand.Mem{Base: secondAddr})

	RET()
	Generate()
}

func loadCtrl16Shuffle(shuffleBase reg.Register, ctrl reg.GPVirtual, upper bool) reg.VecVirtual {
	a := GP64()
	if upper {
		MOVWQZX(ctrl.As16(), a)
		SHRQ(operand.Imm(8), a)
	} else {
		MOVBQZX(ctrl.As8(), a)
	}

	// Left shift by 4 to get the byte level offset for the shuffle table
	SHLQ(operand.Imm(4), a)
	ADDQ(shuffleBase, a)

	shuffle := XMM()
	VLDDQU(operand.Mem{Base: a}, shuffle)
	return shuffle
}

func loadLenValue(ctrl reg.GPVirtual) reg.GPVirtual {
	lt := Load(Param(pLenTable), GP64())
	lenValue := GP64()
	MOVBQZX(ctrl.As8L(), lenValue)
	ADDQ(lt, lenValue)
	return lenValue
}
