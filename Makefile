.PHONY: build
build: 
	go build -v ./cmd/app

.PHONY: run
run: 
	./app.exe
	
.PHONY: bandr
bandr: 
	go build -v ./cmd/app
	./app.exe



# .PHONY: test
# test: 
# 	go test -v -race -timeout 30s ./...
	
.DEFAULT_GOAL := bandr