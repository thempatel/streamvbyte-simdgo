package main

import (
	"bufio"
	"bytes"
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/theMPatel/streamvbyte-simdgo/pkg/util"
)

const (
	goos = "goos"
	goarch = "goarch"
	pkg = "pkg"
	cpu = "cpu"
	benchmark = "Benchmark"
	dashes = "--"
	startSentinel = "## Benchmarks\n\n```text\n"
	endSentinel = "```\n"
)

var (
	validPrefixes = []string{
		goos,
		goarch,
		pkg,
		cpu,
		benchmark,
	}

	readmeFile = filepath.Join(os.Getenv("SBYTE_HOME"), "README.md")
	fWriteOut = flag.Bool("w", false, "write out to readme")
)

func anyPrefix(in string) bool {
	for _, p := range validPrefixes {
		if strings.HasPrefix(in, p) {
			return true
		}
	}

	return false
}

func main() {
	flag.Parse()
	var (
		lines []string
		hasBench = false
	)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if anyPrefix(line) {
			emitDash := strings.HasPrefix(line, cpu)
			hasBench = hasBench || strings.HasPrefix(line, benchmark)
			emitNewline := hasBench && strings.HasPrefix(line, goos)

			if emitNewline {
				lines = append(lines, "")
			}
			if hasBench {
				line = strings.TrimPrefix(line, benchmark)
			}
			lines = append(lines, line)
			if emitDash {
				lines = append(lines, dashes)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("failed to read input: %s", err)
	}

	outputFile, err := os.Open(readmeFile)
	if err != nil {
		log.Fatalf("failed to open file: %s, %s", readmeFile, err)
	}

	allData, err := io.ReadAll(outputFile)
	if err != nil {
		log.Fatalf("failed to read file: %s, %s", readmeFile, err)
	}

	util.SilentClose(outputFile)

	bStart := []byte(startSentinel)
	bEnd := []byte(endSentinel)

	start := bytes.Index(allData, bStart)
	if start < 0 {
		log.Fatalf("couldn't find start sentinel")
	}

	restStart := bytes.Index(allData, bEnd)

	var final []byte
	final = append(final, allData[:start]...)
	final = append(final, bStart...)
	final = append(final, []byte(strings.Join(lines, "\n"))...)
	final = append(final, '\n')
	final = append(final, allData[restStart:]...)

	var out io.Writer
	if *fWriteOut {
		outputFile, err = os.Create(readmeFile)
		if err != nil {
			log.Fatalf("failed to open file: %s, %s", readmeFile, err)
		}
		defer util.SilentClose(outputFile)
		out = outputFile
	} else {
		out = os.Stdout
	}

	_, err = out.Write(final)
	if err != nil {
		log.Fatalf("failed to write to file: %s, %s", readmeFile, err)
	}
}
