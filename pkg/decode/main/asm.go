package main

import (
	"fmt"

	. "github.com/mmcloughlin/avo/build"
	"github.com/mmcloughlin/avo/operand"
	"github.com/theMPatel/streamvbyte-simdgo/pkg/shared"
)

const (
	name = "get8uint32Fast"

	pIn       = "in"
	pOut      = "out"
	pCtrl     = "ctrl"
	pShuffle  = "shuffle"
	pLenTable = "lenTable"
	pR        = "r"
)

var (
	signature = fmt.Sprintf(
		"func(%s []byte, %s []uint32, %s uint16, %s *[256][16]uint8, %s *[256]uint8) (%s uint64)",
		pIn, pOut, pCtrl, pShuffle, pLenTable, pR)
)

func main() {
	TEXT(name, NOSPLIT, signature)

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

	upperAddr, upperSize := shared.LenValueAddr(ctrl, true, pLenTable)
	MOVBQZX(upperAddr, upperSize)
	ADDQ(upperSize, lowerSize)
	Store(lowerSize, Return(pR))

	firstFour := XMM()
	secondFour := XMM()
	VLDDQU(operand.Mem{Base: firstBlock}, firstFour)
	VLDDQU(operand.Mem{Base: secondBlock}, secondFour)

	VPSHUFB(shuffleA, firstFour, firstFour)
	VPSHUFB(shuffleB, secondFour, secondFour)

	outBase := operand.Mem{Base: Load(Param(pOut).Base(), GP64())}

	VMOVDQU(firstFour, outBase)
	VMOVDQU(secondFour, outBase.Offset(16))

	RET()
	Generate()
}
