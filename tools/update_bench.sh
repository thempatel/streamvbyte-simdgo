#!/usr/bin/env bash

set -e

TMP_FILE=$(mktemp)
go test -bench . ./pkg/... | tee $TMP_FILE
cat $TMP_FILE | go run $SBYTE_HOME/tools/parse_and_write_bench.go -w
rm $TMP_FILE