.DEFAULT_GOAL := spec

build:
	go install

ci-publish-release: cross-compile
	go get github.com/tcnksm/ghr
	ghr -u cucumber -token "${GITHUB_TOKEN}" "${CIRCLE_TAG}" dist

cross-compile:
	go get github.com/mitchellh/gox
	gox -ldflags "-X github.com/cucumber/cucumber-engine/src/cli.version=${CIRCLE_TAG}" -output "dist/{{.Dir}}-{{.OS}}-{{.Arch}}"

fix-lint:
	./bin/golangci-lint run --fix -E goimports
	go mod tidy

lint:
	./bin/golangci-lint run -E goimports
	go mod verify

setup:
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s v1.16.0
	go get -u github.com/onsi/ginkgo/ginkgo

spec: build lint unit-test

unit-test:
	ginkgo src/...
