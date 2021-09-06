
generate:
	go generate ./pkg/...

test:
	go test -v ./pkg/...

update-bench:
	./tools/update_bench.sh

fmtgo:
	find ./pkg -type f -iname "*.go" | xargs gofmt -w