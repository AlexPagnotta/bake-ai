package main

// version is set at build time via -ldflags (see .goreleaser.yaml); "dev" otherwise.
var version = "dev"

func main() {
	Execute()
}
