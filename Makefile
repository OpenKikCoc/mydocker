#BUILD_TAGS?=
#BUILD_FLAGS = -ldflags "-X github.com/OpenKikCoc/mydocker/version.GitCommit=`git rev-parse HEAD`"

default: clean build

clean:
	rm -rf bin

build:
	#CGO_ENABLED=1 go build $(BUILD_FLAGS) -o bin/docker ./cmd
	go build $(BUILD_FLAGS) -o bin/docker ./cmd

ubuntu:
	#CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o bin/docker ./cmd
	GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o bin/docker ./cmd

testdir:
	GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o bin/test ./test