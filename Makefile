NAME=coffer
ARCH=$(shell uname -m)
VERSION=2.0.1
ITERATION := 1

default: deps compile

deps:
	go get github.com/c4milo/github-release
	go get github.com/mitchellh/gox
	glide install

compile: deps
	@rm -rf build/
	@gox -ldflags "-X main.Version=$(VERSION)" \
	-osarch="darwin/amd64" \
	-osarch="linux/i386" \
	-osarch="linux/amd64" \
	-osarch="windows/amd64" \
	-osarch="windows/i386" \
	-output "build/{{.Dir}}_$(VERSION)_{{.OS}}_{{.Arch}}/$(NAME)" \
	$(shell glide novendor)

dist:
	$(eval FILES := $(shell ls build))
	@rm -rf dist && mkdir dist
	@for f in $(FILES); do \
		(cd $(shell pwd)/build/$$f && tar -cvzf ../../dist/$$f.tar.gz *); \
		(cd $(shell pwd)/dist && shasum -a 512 $$f.tar.gz > $$f.sha512); \
		echo $$f; \
	done

release:
	@github-release "v$(VERSION)" dist/* --commit "$(git rev-parse HEAD)" --github-repository versent/$(NAME)

test: deps
	go test -cover -v $(shell glide novendor)

clean:
	@rm -rf dist build

.PHONY: default deps clean compile dist release test