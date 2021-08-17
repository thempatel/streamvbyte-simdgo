
generate: generate-deps
	go generate ./internal/...

generate-deps:
	go install internal/shared/main/gentables.go