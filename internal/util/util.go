package util

import "io"

func SilentClose(closer io.Closer) {
	_ = closer.Close()
}
