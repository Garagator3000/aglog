.PHONY: aglog
.DEFAULT_GOAL := aglog

aglog:
	GOOS=linux GOARCH=amd64 go build -ldflags "-w -s" -o bin/aglog ./cmd/aglog

run: aglog
	@./bin/aglog
