# This Makefile is meant to be used by people that do not usually work
# with Go source code. If you know what GOPATH is then you probably
# don't need to bother with make.

.PHONY: gaxis android ios gaxis-cross swarm evm all test clean
.PHONY: gaxis-linux gaxis-linux-386 gaxis-linux-amd64 gaxis-linux-mips64 gaxis-linux-mips64le
.PHONY: gaxis-linux-arm gaxis-linux-arm-5 gaxis-linux-arm-6 gaxis-linux-arm-7 gaxis-linux-arm64
.PHONY: gaxis-darwin gaxis-darwin-386 gaxis-darwin-amd64
.PHONY: gaxis-windows gaxis-windows-386 gaxis-windows-amd64

GOBIN = $(shell pwd)/build/bin
root=$(shell pwd)
GO ?= latest

PKG = ./...

gaxis:
	build/env.sh go run build/ci.go install ./cmd/gaxis
	@echo "Done building."
	@echo "Run \"$(GOBIN)/gaxis\" to launch gaxis."

swarm:
	build/env.sh go run build/ci.go install ./cmd/swarm
	@echo "Done building."
	@echo "Run \"$(GOBIN)/swarm\" to launch swarm."

all:
	build/env.sh go run build/ci.go install

test: all
	build/env.sh go run build/ci.go test

lint: ## Run linters.
	build/env.sh go run build/ci.go lint $(PKG)

clean:
	./build/clean_go_build_cache.sh
	rm -fr build/_workspace/pkg/ $(GOBIN)/*

# The devtools target installs tools required for 'go generate'.
# You need to put $GOBIN (or $GOPATH/bin) in your PATH to use 'go generate'.

devtools:
	env GOBIN= go get -u golang.org/x/tools/cmd/stringer
	env GOBIN= go get -u github.com/kevinburke/go-bindata/go-bindata
	env GOBIN= go get -u github.com/fjl/gencodec
	env GOBIN= go get -u github.com/golang/protobuf/protoc-gen-go
	env GOBIN= go install ./cmd/abigen
	@type "npm" 2> /dev/null || echo 'Please install node.js and npm'
	@type "solc" 2> /dev/null || echo 'Please install solc'
	@type "protoc" 2> /dev/null || echo 'Please install protoc'

# Cross Compilation Targets (xgo)

gaxis-cross: gaxis-linux gaxis-darwin gaxis-windows
	@echo "Full cross compilation done:"
	@ls -ld $(GOBIN)/gaxis-*

gaxis-linux: gaxis-linux-amd640-v3 gaxis-linux-amd64-v4
	@echo "Linux cross compilation done:"
	@ls -ld $(GOBIN)/gaxis-linux-*

gaxis-linux-amd64-v3:
	build/env.sh linux-v3 go run build/ci.go xgo -- --go=$(GO) --out=gaxis-v3 --targets=linux/amd64 -v ./cmd/gaxis
	#build/env.sh linux-v3 go run build/ci.go xgo -- --go=$(GO) --out=bootnode-v3 --targets=linux/amd64 -v ./cmd/bootnode
	@echo "Linux centos amd64 cross compilation done:"
	@ls -ld $(GOBIN)/gaxis-v3-linux-* | grep amd64

gaxis-linux-amd64-v4:
	build/env.sh linux-v4 go run build/ci.go xgo -- --go=$(GO) --out=gaxis-v4 --targets=linux/amd64 -v ./cmd/gaxis
	#build/env.sh linux-v3 go run build/ci.go xgo -- --go=$(GO) --out=bootnode-v4 --targets=linux/amd64 -v ./cmd/bootnode
	@echo "Linux  ubuntu amd64 cross compilation done:"
	@ls -ld $(GOBIN)/gaxis-v4-linux-* | grep amd64

gaxis-darwin: gaxis-darwin-amd64
	@echo "Darwin cross compilation done:"
	@ls -ld $(GOBIN)/gaxis-darwin-*


gaxis-darwin-amd64:
	build/env.sh darwin-amd64 go run build/ci.go xgo -- --go=$(GO) --targets=darwin/amd64 -v ./cmd/gaxis
	@echo "Darwin amd64 cross compilation done:"
	@ls -ld $(GOBIN)/gaxis-darwin-* | grep amd64

gaxis-windows: gaxis-windows-amd64
	@echo "Windows cross compilation done:"
	@ls -ld $(GOBIN)/gaxis-windows-*

gaxis-windows-amd64:
	build/env.sh windows-amd64 go run build/ci.go xgo -- --go=$(GO)  --targets=windows/amd64 -v ./cmd/gaxis
	@echo "Windows amd64 cross compilation done:"
	@ls -ld $(GOBIN)/gaxis-windows-* | grep amd64

gaxistx-darwin-amd64:
	build/env.sh darwin-amd64 go run build/ci.go xgo -- --go=$(GO) --out=gaxistx  --targets=darwin/amd64 -v ./cmd/tx
	@echo "Darwin amd64 cross compilation done:"
	@ls -ld $(GOBIN)/gaxistx-darwin-* | grep amd64

gaxistx-linux-amd64-v3:
	build/env.sh linux-v3 go run build/ci.go xgo -- --go=$(GO) --out=gaxistx-v3 --targets=linux/amd64 -v ./cmd/tx
	@echo "Linux centos amd64 cross compilation done:"
	@ls -ld $(GOBIN)/gaxistx-v3-linux-* | grep amd64

gaxistx-linux-amd64-v4:
	build/env.sh linux-v4 go run build/ci.go xgo -- --go=$(GO) --out=gaxistx-v4 --targets=linux/amd64 -v ./cmd/tx
	@echo "Linux  ubuntu amd64 cross compilation done:"
	@ls -ld $(GOBIN)/gaxistx-v4-linux-* | grep amd64

gaxistx-windows-amd64:
	build/env.sh windows-amd64 go run build/ci.go xgo -- --go=$(GO) --out=gaxistx --targets=windows/amd64 -v ./cmd/tx
	@echo "Windows amd64 cross compilation done:"
	@ls -ld $(GOBIN)/gaxistx-windows-* | grep amd64
