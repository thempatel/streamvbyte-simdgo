
generate: generate-deps
	go generate ./shared/...

generate-deps:
	go install shared/main/gentables.go