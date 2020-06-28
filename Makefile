#BUILD_TAGS?=
#BUILD_FLAGS = -ldflags "-X github.com/OpenKikCoc/mydocker/version.GitCommit=`git rev-parse HEAD`"

default: clean build

clean:
	rm -rf bin

build:
	go build $(BUILD_FLAGS) -o bin/docker ./cmd

ubuntu:
	GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o bin/docker ./cmd

testdir:
	GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o bin/test ./test