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
	name      = "Get8uint32FastAsm"
	nameDelta = "Get8uint32DeltaFastAsm"

	pIn       = "in"
	pOut      = "out"
	pCtrl     = "ctrl"
	pShuffle  = "shuffle"
	pLenTable = "lenTable"
	pPrev     = "prev"
)

var (
	signature = fmt.Sprintf(
		"func(%s []byte, %s []uint32, %s uint16, %s *[256][16]uint8, %s *[256]uint8)",
		pIn, pOut, pCtrl, pShuffle, pLenTable)

	signatureDelta = fmt.Sprintf(
		"func(%s []byte, %s []uint32, %s uint16, %s uint32, %s *[256][16]uint8, %s *[256]uint8)",
		pIn, pOut, pCtrl, pPrev, pShuffle, pLenTable)
)

func main() {
	regular()
	differential()
	Generate()
}

func regular() {
	TEXT(name, NOSPLIT, signature)

	firstFour, secondFour := coreAlgorithm()
	outBase := operand.Mem{Base: Load(Param(pOut).Base(), GP64())}

	VMOVDQU(firstFour, outBase)
	VMOVDQU(secondFour, outBase.Offset(16))

	RET()
}

func differential() {
	TEXT(nameDelta, NOSPLIT, signatureDelta)

	firstFour, secondFour := coreAlgorithm() // [A B C D] [E F G H]
	prevSingular, err := Param(pPrev).Resolve()
	if err != nil {
		log.Fatalf("failed to get addr of prev")
	}

	prev := XMM()
	VBROADCASTSS(prevSingular.Addr, prev) // [P P P P]
	undoDelta(firstFour, prev)

	VPSHUFD(operand.Imm(0xff), firstFour, prev) // [A B C D] -> [D D D D]
	undoDelta(secondFour, prev)

	outBase := operand.Mem{Base: Load(Param(pOut).Base(), GP64())}

	VMOVDQU(firstFour, outBase)
	VMOVDQU(secondFour, outBase.Offset(16))

	RET()
}

func undoDelta(four, prev reg.VecVirtual) {
	adder := XMM()                       // [A B C D]
	VPSLLDQ(operand.Imm(4), four, adder) // [- A  B  C]
	VPADDD(four, adder, four)            // [A AB BC CD]
	VPSLLDQ(operand.Imm(8), four, adder) // [- - A AB]
	VPADDD(four, prev, four)             // [PA PAB PBC PCD]
	VPADDD(four, adder, four)            // [PA PAB PABC PABCD]
}

func coreAlgorithm() (reg.VecVirtual, reg.VecVirtual) {
	ctrl := GP64()
	Load(Param(pCtrl), ctrl)

	shuffleBase := Load(Param(pShuffle), GP64())
	shuffleA := shared.CalculateShuffleAddrFromCtrl(shuffleBase, ctrl, false)
	shuffleB := shared.CalculateShuffleAddrFromCtrl(shuffleBase, ctrl, true)

	firstBlock := Load(Param(pIn).Base(), GP64())
	secondBlock := GP64()
	MOVQ(firstBlock, secondBlock)
	lowerAddr, lowerSize := shared.LenValueAddr(ctrl, false, pLenTable)

	MOVBQZX(lowerAddr, lowerSize)
	ADDQ(lowerSize, secondBlock)

	firstFour := XMM()
	secondFour := XMM()
	VLDDQU(operand.Mem{Base: firstBlock}, firstFour)
	VLDDQU(operand.Mem{Base: secondBlock}, secondFour)

	VPSHUFB(shuffleA, firstFour, firstFour)
	VPSHUFB(shuffleB, secondFour, secondFour)

	return firstFour, secondFour
}
