.DEFAULT_GOAL := spec

build:
	go install

fix-lint:
	goimports -w src

lint:
	goimports -l src
	gometalinter.v2

setup:
	go get -u \
		github.com/Masterminds/glide \
		gopkg.in/alecthomas/gometalinter.v2 \
		github.com/onsi/ginkgo/ginkgo
	glide install
	gometalinter.v2 --install

spec: build lint unit-test

unit-test:
	ginkgo src/...
