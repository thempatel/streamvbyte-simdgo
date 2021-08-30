package main

import (
	"fmt"

	. "github.com/mmcloughlin/avo/build"
	"github.com/mmcloughlin/avo/operand"
	"github.com/mmcloughlin/avo/reg"
)

const (
	name = "get8uint32Fast"

	pIn    = "in"
	pOut   = "out"
	pCtrl  = "ctrl"
	pShuffle = "shuffle"
	pLenTable  = "lenTable"
	pR 			= "r"
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
	shuffleA := loadCtrl16Shuffle(shuffleBase, ctrl, false)
	shuffleB := loadCtrl16Shuffle(shuffleBase, ctrl, true)

	firstBlock := Load(Param(pIn).Base(), GP64())
	secondBlock := GP64()
	MOVQ(firstBlock, secondBlock)
	lowerSize := loadLenValue(ctrl, false)

	MOVBQZX(operand.Mem{Base: lowerSize}, lowerSize)
	ADDQ(lowerSize, secondBlock)

	upperSize := loadLenValue(ctrl, true)
	MOVBQZX(operand.Mem{Base: upperSize}, upperSize)
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

func loadCtrl16Shuffle(shuffleBase reg.Register, ctrl reg.GPVirtual, upper bool) operand.Mem {
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

	return operand.Mem{Base: a}
}

func loadLenValue(ctrl reg.GPVirtual, upper bool) reg.GPVirtual {
	lt := Load(Param(pLenTable), GP64())
	lenValue := GP64()
	if upper {
		MOVWQZX(ctrl.As16(), lenValue)
		SHRQ(operand.Imm(8), lenValue)
	} else {
		MOVBQZX(ctrl.As8L(), lenValue)
	}
	ADDQ(lt, lenValue)
	return lenValue
}