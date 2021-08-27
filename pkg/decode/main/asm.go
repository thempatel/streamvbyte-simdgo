package main

import (
	"fmt"

	. "github.com/mmcloughlin/avo/build"
	"github.com/mmcloughlin/avo/operand"
)

const (
	name = "get8uint32Fast"

	pIn = "in"
	pOut = "out"
	pShufA = "shufA"
	pShufB = "shufB"
	pLenA = "lenA"
)

var (
	signature = fmt.Sprintf(
		"func(%s []byte, %s []uint32, %s, %s *[16]uint8, %s uint8)",
		pIn, pOut, pShufA, pShufB, pLenA)
)

func main() {
	TEXT(name, NOSPLIT, signature)
	Doc("get8uint32Fast decodes 8 32-bit unsigned integers at a time.")

	shuffleA := operand.Mem{
		Base: Load(Param(pShufA), GP64()),
	}

	shuffleB := operand.Mem{
		Base: Load(Param(pShufB), GP64()),
	}

	firstBlock := Load(Param(pIn).Base(), GP64())
	secondBlock := GP64()
	MOVQ(firstBlock, secondBlock)
	size := Load(Param(pLenA), GP64())
	ADDQ(size, secondBlock)

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