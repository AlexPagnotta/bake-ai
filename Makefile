.PHONY: dev build test snapshot release

dev:                ## run the CLI locally
	go run ./cmd/bake

build:              ## build everything
	go build ./...

test:               ## run tests
	go test ./...

snapshot:           ## test a release build locally (no publish)
	goreleaser build --clean --snapshot

release:            ## publish a release (tag must already be pushed — see the cut-release skill)
	goreleaser release --clean
