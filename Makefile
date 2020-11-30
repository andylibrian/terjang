.PHONY: default
default:
	mkdir -p ./bin
	go build -o ./bin/terjang ./cmd/terjang/
