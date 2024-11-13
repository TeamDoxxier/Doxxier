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
	go build  

