client:
	@rm -f client.*
	@go build ./client/
	@./client.exe
server:
	@rm -f server.*
	@go build ./server/
	@./server.exe

.PHONY: client server