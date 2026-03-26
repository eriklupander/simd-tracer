.PHONY: fmt vet run test

all: fmt vet test

fmt:
	GOEXPERIMENT=simd go fmt ./...

vet:
	GOEXPERIMENT=simd go vet ./...

run:
	GOEXPERIMENT=simd go run main.go

test:
	GOEXPERIMENT=simd go test ./...
