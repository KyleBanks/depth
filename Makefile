RELEASE_PKG = ./cmd/depth

# Builds and installs the depth CLI.
install:
	@go install -v github.com/KyleBanks/depth/cmd/depth
	@echo "depth installed."
.PHONY: install

# Runs a number of depth commands as examples of what's possible.
example: | install
	depth github.com/KyleBanks/depth/cmd/depth strings

	depth -internal strings 

	depth -json github.com/KyleBanks/depth/cmd/depth

	depth -test github.com/KyleBanks/depth/cmd/depth

	depth -test -internal strings

	depth -test -internal -max 3 strings
.PHONY: example

include github.com/KyleBanks/make/git/precommit
include github.com/KyleBanks/make/go/sanity
include github.com/KyleBanks/make/go/release
include github.com/KyleBanks/make/go/bench
