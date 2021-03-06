package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"io"
	"log"
	"os"

	"github.com/pkg/errors"
	"github.com/theMPatel/streamvbyte-simdgo/pkg/util"
)

var (
	fOut     = flag.String("out", "", "path to output")
	fPackage = flag.String("package", "shared", "package name")
)

const MaxControlByte = 1 << 8

func main() {
	flag.Parse()

	if *fOut == "" {
		log.Fatalf("outfile cannot be empty")
	}

	out := &bytes.Buffer{}

	_, _ = fmt.Fprintln(out, "// Code generated by gentables. DO NOT EDIT.")
	_, _ = fmt.Fprintf(out, "\npackage %s\n", *fPackage)

	if err := genPerNumLengthTable(out); err != nil {
		log.Fatalf("failed to gen per num length table")
	}

	if err := genPerQuadLengthTable(out); err != nil {
		log.Fatalf("failed to gen sum length table")
	}

	if err := genEncodeShuffleTable(out); err != nil {
		log.Fatalf("failed to gen encode shuffle table")
	}

	if err := genDecodeShuffleTable(out); err != nil {
		log.Fatalf("failed to gen decode shuffle table")
	}

	final, err := format.Source(out.Bytes())
	if err != nil {
		log.Fatalf("failed to go fmt output")
	}

	fileOut, err := os.Create(*fOut)
	if err != nil {
		log.Fatalf("failed to open: %s %s", *fOut, err)
	}
	defer util.SilentClose(fileOut)

	_, err = fileOut.Write(final)
	if err != nil {
		log.Fatalf("failed to write generated tables to file")
	}
}

func newLineAfter(countPerLine int) func(out io.Writer) {
	count := 1
	return func(out io.Writer) {
		if count%countPerLine == 0 {
			_, _ = fmt.Fprintln(out, "")
		} else {
			_, _ = fmt.Fprintf(out, " ")
		}
		count++
	}
}

func genPerNumLengthTable(out io.Writer) error {
	_, _ = fmt.Fprintf(out, "\nvar PerNumLenTable *[256][4]uint8 = &[256][4]uint8{\n")
	tabber := newLineAfter(4)
	for i := 0; i < MaxControlByte; i++ {
		one, two, three, four := sizes(uint8(i))
		_, err := fmt.Fprintf(out, "\t{%d, %d, %d, %d},", one, two, three, four)
		if err != nil {
			return errors.Wrapf(err, "failed to write per num len: %d", i)
		}
		tabber(out)
	}
	_, _ = fmt.Fprintln(out, "}")
	return nil
}

func genPerQuadLengthTable(out io.Writer) error {
	_, _ = fmt.Fprintf(out, "\nvar PerControlLenTable *[256]uint8 = &[256]uint8{\n")
	tabber := newLineAfter(8)
	for i := 0; i < MaxControlByte; i++ {
		one, two, three, four := sizes(uint8(i))
		_, err := fmt.Fprintf(out, "\t%d,", one+two+three+four)
		if err != nil {
			return errors.Wrapf(err, "failed to write summed len: %d", i)
		}
		tabber(out)
	}
	_, _ = fmt.Fprintln(out, "}")
	return nil
}

const (
	shuffleFmtStr = "%#02x, %#02x, %#02x, %#02x, %#02x, %#02x, %#02x, %#02x, %#02x, %#02x, %#02x, %#02x, %#02x, %#02x, %#02x, %#02x},"
	commentStr    = "\t// %d\t%#02x\t%08b\tlen\t%d\t%d\t%d\t%d\n"
)

func genEncodeShuffleTable(out io.Writer) error {
	_, _ = fmt.Fprintf(out, "\nvar EncodeShuffleTable *[256][16]uint8 = &[256][16]uint8{\n")
	tabber := newLineAfter(1)
	for i := 0; i < MaxControlByte; i++ {
		one, two, three, four := sizes(uint8(i))
		_, _ = fmt.Fprintf(out, commentStr, i, i, i, one, two, three, four)
		_, err := fmt.Fprintf(out, "\t{")
		if err != nil {
			return errors.Wrapf(err, "failed to write encode shuffle table")
		}

		var positions []interface{}
		var base uint8
		for _, size := range []uint8{one, two, three, four} {
			for j := uint8(0); j < size; j++ {
				positions = append(positions, base+j)
			}
			base += 4
		}

		for len(positions) < 16 {
			positions = append(positions, 0xff)
		}
		_, err = fmt.Fprintf(out, shuffleFmtStr, positions...)
		if err != nil {
			return errors.Wrapf(err, "failed to write per num len: %d", i)
		}
		tabber(out)
	}
	_, _ = fmt.Fprintln(out, "}")
	return nil
}

func genDecodeShuffleTable(out io.Writer) error {
	_, _ = fmt.Fprintf(out, "\nvar DecodeShuffleTable *[256][16]uint8 = &[256][16]uint8{\n")
	tabber := newLineAfter(1)
	for i := 0; i < MaxControlByte; i++ {
		one, two, three, four := sizes(uint8(i))
		_, _ = fmt.Fprintf(out, commentStr, i, i, i, one, two, three, four)
		_, err := fmt.Fprintf(out, "\t{")
		if err != nil {
			return errors.Wrapf(err, "failed to write encode shuffle table")
		}

		var positions []interface{}
		var pos uint8
		for _, size := range []uint8{one, two, three, four} {
			for j := 0; j < 4; j++ {
				if size > 0 {
					positions = append(positions, pos)
					pos++
					size--
				} else {
					positions = append(positions, 0xff)
				}
			}
		}

		_, err = fmt.Fprintf(out, shuffleFmtStr, positions...)
		if err != nil {
			return errors.Wrapf(err, "failed to write per num len: %d", i)
		}
		tabber(out)
	}
	_, _ = fmt.Fprintln(out, "}")
	return nil
}

// sizes returns the length in bytes for each of the four numbers
// represented by the provided control byte.
func sizes(control uint8) (one uint8, two uint8, three uint8, four uint8) {
	one = (control & 3) + 1
	two = (control >> 2 & 3) + 1
	three = (control >> 4 & 3) + 1
	four = (control >> 6 & 3) + 1
	return
}
