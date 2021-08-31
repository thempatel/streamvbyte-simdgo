package main

import (
	"fmt"
	"log"

	. "github.com/mmcloughlin/avo/build"
	"github.com/mmcloughlin/avo/operand"
	"github.com/mmcloughlin/avo/reg"
	"github.com/theMPatel/streamvbyte-simdgo/pkg/shared"
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

	firstFour, secondFour := shared.Load8(pIn)
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
	coreAlgorithm(shared.Load8(pIn))
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
	firstShuffle := shared.CalculateShuffleAddrFromCtrl(shuffleBase, ctrl, false)
	secondShuffle := shared.CalculateShuffleAddrFromCtrl(shuffleBase, ctrl, true)

	VPSHUFB(firstShuffle, firstFour, firstFour)
	VPSHUFB(secondShuffle, secondFour, secondFour)

	firstAddr := Load(Param(pOut).Base(), GP64())
	secondAddr := GP64()
	MOVQ(firstAddr, secondAddr)

	lenAddr, lenValue := shared.LenValueAddr(ctrl, false, pLenTable)

	MOVBQZX(lenAddr, lenValue)
	ADDQ(lenValue, secondAddr)

	VMOVDQU(firstFour, operand.Mem{Base: firstAddr})
	VMOVDQU(secondFour, operand.Mem{Base: secondAddr})

	RET()
}
