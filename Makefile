GOTEST_FLAGS=-cpu=1,2,4

BASE_PACKAGE=github.com/stormcat24/circle-warp

deps:
		go get github.com/tools/godep
		godep restore

deps-save:
		godep save $(BASE_PACKAGE)/...

build:
		godep go build -o bin/circle-warp main.go
