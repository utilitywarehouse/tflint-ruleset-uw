TFLINT_PLUGIN_DIR ?= ~/.tflint.d/plugins

default: build

test:
	go test ./...

build:
	go build

install: build
	mkdir -p $(TFLINT_PLUGIN_DIR)
	mv ./tflint-ruleset-uw $(TFLINT_PLUGIN_DIR)
