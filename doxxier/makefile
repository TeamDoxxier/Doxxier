# Makefile for Doxxier
.PHONY: build clean update binary

clean:
	go clean -modcache
	go mod tidy 
	go mod vendor -e 
	rm -rf *.wasm 
	rm -rf assets/wasm/*

update:
	-GOFLAGS="" go get all 

build: 
	GOOS=js GOARCH=wasm go build  

binary:
	mkdir -p assets/wasm
	GOOS=js GOARCH=wasm go build -ldflags '-w -s' -trimpath -o doxxier.wasm main.go
	cp doxxier.wasm assets/wasm/
	cp doxxier.wasm ../test-harness/static/assets/
	rm -rf doxxier.wasm
