package shared

import (
	. "github.com/mmcloughlin/avo/build"
	"github.com/mmcloughlin/avo/operand"
	"github.com/mmcloughlin/avo/reg"
)

func CalculateShuffleAddrFromCtrl(shuffleBase reg.Register, ctrl reg.GPVirtual, upper bool) operand.Mem {
	addr := GP64()
	if upper {
		MOVWQZX(ctrl.As16(), addr)
		SHRQ(operand.Imm(8), addr)
	} else {
		MOVBQZX(ctrl.As8(), addr)
	}

	// Left shift by 4 to get the byte level offset for the shuffle table
	SHLQ(operand.Imm(4), addr)
	ADDQ(shuffleBase, addr)

	return operand.Mem{Base: addr}
}

func LenValueAddr(ctrl reg.GPVirtual, upper bool, lenTableParam string) (operand.Mem, reg.GPVirtual) {
	lenTableBase := Load(Param(lenTableParam), GP64())
	lenValueAddr := GP64()
	if upper {
		MOVWQZX(ctrl.As16(), lenValueAddr)
		SHRQ(operand.Imm(8), lenValueAddr)
	} else {
		MOVBQZX(ctrl.As8L(), lenValueAddr)
	}
	ADDQ(lenTableBase, lenValueAddr)

	return operand.Mem{Base: lenValueAddr}, lenValueAddr
}

func Load8(paramName string) (reg.VecVirtual, reg.VecVirtual) {
	arrBase := operand.Mem{
		Base: Load(Param(paramName).Base(), GP64()),
	}
	firstFour := XMM()
	secondFour := XMM()
	VLDDQU(arrBase, firstFour)
	VLDDQU(arrBase.Offset(16), secondFour)

	return firstFour, secondFour
}
