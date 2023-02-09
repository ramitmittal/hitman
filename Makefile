.DEFAULT_GOAL = build

GIT_SHA := $(shell git rev-parse --short HEAD)
BINARY_NAME=out/hitman

download:
	go mod download

generate: download
	go generate ./...

test:
	go test -shuffle=on -v ./...

clean:
	go clean
	rm -rf out/

build: generate
	GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME}-linux -ldflags "-X main.GitSHA=$(GIT_SHA)" cmd/hitman/*
	GOARCH=amd64 GOOS=darwin go build -o ${BINARY_NAME}-darwin -ldflags "-X main.GitSHA=$(GIT_SHA)" cmd/hitman/*
	GOARCH=amd64 GOOS=windows go build -o ${BINARY_NAME}-windows.exe -ldflags "-X main.GitSHA=$(GIT_SHA)" cmd/hitman/*

run: build
	./${BINARY_NAME}-linux
