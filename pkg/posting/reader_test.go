package posting

import (
	"fmt"
	"testing"
)

func TestWhatHappens(t *testing.T) {
	count := 14
	fmt.Printf("%#b\n", -8)
	fmt.Println(count &^7)
}
