.PHONY: build
build: 
	go build -v .\server.go


.PHONY: bandr
go:
	go build -v .\server.go
	./server.exe



# .PHONY: test
# test: 
# 	go test -v -race -timeout 30s ./...
	
.DEFAULT_GOAL := go