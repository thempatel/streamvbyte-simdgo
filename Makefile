
generate:
	go generate ./pkg/...

test:
	go test -v ./pkg/...

update-bench:
	./tools/update_bench.sh