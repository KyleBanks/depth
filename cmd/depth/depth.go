package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/KyleBanks/depth"
)

const (
	outputPadding    = "  "
	outputPrefix     = "├ "
	outputPrefixLast = "└ "
)

var outputJSON bool

func main() {
	var t depth.Tree
	flag.BoolVar(&t.ResolveInternal, "internal", false, "If set, resolves dependencies of internal (stdlib) packages.")
	flag.BoolVar(&t.ResolveTest, "test", false, "If set, resolves dependencies used for testing.")
	flag.IntVar(&t.MaxDepth, "max", 0, "Sets the maximum depth of dependencies to resolve.")
	flag.BoolVar(&outputJSON, "json", false, "If set, outputs the depencies in JSON format.")
	flag.Parse()

	if err := handlePkgs(&t, flag.Args()); err != nil {
		os.Exit(1)
	}
}

// handlePkgs takes a slice of package names, resolves a Tree on them,
// and outputs each Tree to Stdout.
func handlePkgs(t *depth.Tree, pkgs []string) error {
	for _, pkg := range pkgs {
		err := t.Resolve(pkg)
		if err != nil {
			fmt.Printf("'%v': FATAL: %v\n", pkg, err)
			return err
		}

		if outputJSON {
			writePkgJSON(os.Stdout, *t.Root)
			continue
		}

		writePkg(os.Stdout, *t.Root, 0, false)
	}

	return nil
}

// writePkgJSON writes the full Pkg as JSON to the provided Writer.
func writePkgJSON(w io.Writer, p depth.Pkg) {
	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	e.Encode(p)
}

// writePkg recursively prints a Pkg and its dependencies to the Writer provided.
func writePkg(w io.Writer, p depth.Pkg, indent int, isLast bool) {
	var prefix string
	if indent > 0 {
		prefix = outputPrefix

		if isLast {
			prefix = outputPrefixLast
		}
	}

	out := fmt.Sprintf("%v%v%v\n", strings.Repeat(outputPadding, indent), prefix, p.String())
	w.Write([]byte(out))

	for idx, d := range p.Deps {
		writePkg(w, d, indent+1, idx == len(p.Deps)-1)
	}
}
