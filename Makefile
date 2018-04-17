.DEFAULT_GOAL := spec

build:
	go install

ci-publish-release: cross-compile
	go get github.com/tcnksm/ghr
	ghr -u cucumber -token "${GITHUB_TOKEN}" "${CIRCLE_TAG}" dist

cross-compile:
	go get github.com/mitchellh/gox
	gox -ldflags "-X github.com/cucumber/cucumber-pickle-runner/src/cli.version=${CIRCLE_TAG}" -output "dist/{{.Dir}}-{{.OS}}-{{.Arch}}"

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
