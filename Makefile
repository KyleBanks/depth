include github.com/KyleBanks/make/go/sanity
include github.com/KyleBanks/make/go/release

# Builds and installs the depth CLI.
install:
	@go install github.com/KyleBanks/depth/cmd/depth
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
