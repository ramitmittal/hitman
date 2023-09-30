.DEFAULT_GOAL = build

GIT_SHA := $(shell git rev-parse --short HEAD)
GIT_TAG := $(shell git describe --tags --exact-match --always)
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

build: generate clean
	GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME}_${GIT_TAG}_linux_amd64 -ldflags "-X main.GitSHA=$(GIT_SHA) -X main.GitTag=$(GIT_TAG)" cmd/hitman/*
	GOARCH=amd64 GOOS=darwin go build -o ${BINARY_NAME}_${GIT_TAG}_darwin_amd64 -ldflags "-X main.GitSHA=$(GIT_SHA) -X main.GitTag=$(GIT_TAG)" cmd/hitman/*
	GOARCH=amd64 GOOS=windows go build -o ${BINARY_NAME}_${GIT_TAG}_windows_amd64.exe -ldflags "-X main.GitSHA=$(GIT_SHA) -X main.GitTag=$(GIT_TAG)" cmd/hitman/*