package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/theMPatel/streamvbyte-simdgo/pkg/shared"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const (
	lowerCtrlMask uint16 = 0xff
	upperCtrlMask uint16 = ^lowerCtrlMask
)

func doSomething() uint8 {
	var control = uint16(rand.Int31n(256))
	lowerIdx := control & lowerCtrlMask
	res := shared.PerControlLenTable[lowerIdx]
	return res+1
}

func main() {
	fmt.Println(doSomething())
}
