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

# Runs test suit, vet, golint, and fmt.
sanity:
	@echo "---------------- TEST ----------------"
	@go list ./... | grep -v vendor/ | xargs go test --cover 

	@echo "---------------- VET ----------------"
	@go list ./... | grep -v vendor/ | xargs go vet 

	@echo "---------------- LINT ----------------"
	@go list ./... | grep -v vendor/ | xargs golint

	@echo "---------------- FMT ----------------"
	@go list ./... | grep -v vendor/ | xargs go fmt
.PHONY: sanity

# Creates release binaries for each supported platform/architecture.
release: | sanity
	@gox -osarch="darwin/386 darwin/amd64 linux/386 linux/amd64 linux/arm windows/386 windows/amd64" \
		-output "bin/{{.Dir}}_{{.OS}}_{{.Arch}}" .
.PHONY: release