package shared

//go:generate go run ./main/gentables.go -out ./tables.go -package shared

func ControlByteToSize(in uint8) int {
	return int(PerControlLenTable[in])
}

func ControlByteToSizeTwo(in uint16) int {
	return int(PerControlLenTable[in&0xff] + PerControlLenTable[in>>8])
}
