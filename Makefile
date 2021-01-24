.PHONY: default
default:
	mkdir -p ./bin
	go build -o ./bin/terjang ./cmd/terjang/

prepare:
	cd web; npm ci; cd -
	go get github.com/rakyll/statik

build-ui:
	cd web; ./node_modules/.bin/vue-cli-service build; cd -
	statik -src=./web/dist -dest=./pkg/server
