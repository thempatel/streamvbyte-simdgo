package main

import (
	"fmt"
	"log"

	. "github.com/mmcloughlin/avo/build"
	"github.com/mmcloughlin/avo/operand"
	"github.com/mmcloughlin/avo/reg"
)

const (
	name      = "put8uint32Fast"
	nameDiff  = "put8uint32DiffFast"
	pIn       = "in"
	pOut      = "outBytes"
	pShuffle  = "shuffle"
	pLenTable = "lenTable"
	pPrev     = "prev"
	pR        = "r"
)

var (
	signature = fmt.Sprintf(
		"func(%s []uint32, %s []byte, %s *[256][16]uint8, %s *[256]uint8) (%s uint16)",
		pIn, pOut, pShuffle, pLenTable, pR)

	signatureDiff = fmt.Sprintf(
		"func(%s []uint32, %s []byte, %s uint32, %s *[256][16]uint8, %s *[256]uint8) (%s uint16)",
		pIn, pOut, pPrev, pShuffle, pLenTable, pR)

	mask1111R = ConstData("mask0101", operand.U16(0x0101))
	mask7F00R = ConstData("mask7F00", operand.U16(0x7F00))
)

func main() {
	regular()
	differential()
	Generate()
}

func differential() {
	TEXT(nameDiff, NOSPLIT, signatureDiff)

	prevSingular, err := Param(pPrev).Resolve()
	if err != nil {
		log.Fatalf("failed to get addr of prev")
	}

	firstFour, secondFour := load8(pIn)
	prev := XMM()
	VPALIGNR(operand.Imm(12), firstFour, secondFour, prev)
	VPSUBD(prev, secondFour, secondFour)

	VBROADCASTSS(prevSingular.Addr, prev)
	VPALIGNR(operand.Imm(12), prev, firstFour, prev)
	VPSUBD(prev, firstFour, firstFour)

	coreAlgorithm(firstFour, secondFour)
}

func regular() {
	TEXT(name, NOSPLIT, signature)
	coreAlgorithm(load8(pIn))
}

func coreAlgorithm(firstFour, secondFour reg.VecVirtual) {
	onesMask := XMM()
	sevenFzerozero := XMM()
	VPBROADCASTW(mask1111R, onesMask)
	VPBROADCASTW(mask7F00R, sevenFzerozero)

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

	lenValue := loadLenValue(pLenTable, ctrl)

	MOVBQZX(operand.Mem{Base: lenValue}, lenValue)
	ADDQ(lenValue, secondAddr)

	VMOVDQU(firstFour, operand.Mem{Base: firstAddr})
	VMOVDQU(secondFour, operand.Mem{Base: secondAddr})

	RET()
}

func load8(paramName string) (reg.VecVirtual, reg.VecVirtual) {
	arrBase := operand.Mem{
		Base: Load(Param(paramName).Base(), GP64()),
	}
	firstFour := XMM()
	secondFour := XMM()
	VLDDQU(arrBase, firstFour)
	VLDDQU(arrBase.Offset(16), secondFour)

	return firstFour, secondFour
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

func loadLenValue(paramName string, ctrl reg.GPVirtual) reg.GPVirtual {
	lt := Load(Param(paramName), GP64())
	lenValue := GP64()
	MOVBQZX(ctrl.As8L(), lenValue)
	ADDQ(lt, lenValue)
	return lenValue
}
