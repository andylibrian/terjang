.PHONY: default
default:
	mkdir -p ./bin
	CGO_ENABLED=0 go build -a -o ./bin/terjang ./cmd/terjang/

prepare:
	cd web; npm ci; cd -
	go get github.com/rakyll/statik

build-ui:
	cd web; ./node_modules/.bin/vue-cli-service build; cd -
	statik -src=./web/dist -dest=./pkg/server
