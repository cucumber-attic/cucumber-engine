# Contributing

Before anything else, thank you. Thank you for taking some of your precious time helping this project move forward.

## Development

#### Requirements

* [golang](https://golang.org/)
  * Clone this repo into `$GOPATH/src/github.com/cucumber/cucumber-engine`. See [here](https://golang.org/doc/code.html#Organization) for docs.
  * Add `$GOPATH/bin` to your `$PATH`. 
* [make](https://www.gnu.org/software/make/)

#### Setup

* `make setup`

#### Run linter / tests

* `make spec`

## Release Process

* Ensure the `CHANGELOG.md` is up to date
* Tag the commit
* A release will be created on CI
