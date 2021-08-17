
generate: generate-deps
	go generate ./...

generate-deps:
	go install shared/main/gentables.go