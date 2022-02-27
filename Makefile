build:
	go build -v .\server.go

go:
	go build -v .\server.go
	./server.exe


deploy:
	git rebase master deploy
	git status
	git push
	git switch master

.DEFAULT_GOAL := go