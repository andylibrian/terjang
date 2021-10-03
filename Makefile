.PHONY: default
default:
	mkdir -p ./bin
	CGO_ENABLED=0 go build -a -o ./bin/terjang ./cmd/terjang/

prepare:
	cd web; npm ci; cd -

build-ui:
	cd web; ./node_modules/.bin/vue-cli-service build; cd -
