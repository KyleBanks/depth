package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/KyleBanks/depth"
)

const (
	outputPadding    = "  "
	outputPrefix     = "├ "
	outputPrefixLast = "└ "
)

var outputJSON bool

func main() {
	ResolveInternal := flag.Bool("internal", false, "If set, resolves dependencies of internal (stdlib) packages.")
	ResolveTest := flag.Bool("test", false, "If set, resolves dependencies used for testing.")
	MaxDepth := flag.Int("max", 0, "Sets the maximum depth of dependencies to resolve.")
	outputJSON := flag.Bool("json", false, "If set, outputs the depencies in JSON format.")
	flag.Parse()

	pkgs := flag.Args()
	t := &depth.Tree{
		ResolveInternal: *ResolveInternal,
		ResolveTest:     *ResolveTest,
		MaxDepth:        *MaxDepth,
	}

	handlePkgs(t, pkgs, *outputJSON)

}

func handlePkgs(t *depth.Tree, pkgs []string, outputJSON bool) {
	var wg sync.WaitGroup
	for _, pkg := range pkgs {
		wg.Add(1)
		go handlePkg(&wg, t, pkg, outputJSON)
	}
	wg.Wait()
}

func handlePkg(wg *sync.WaitGroup, t *depth.Tree, pkg string, outputJSON bool) {
	defer wg.Done()
	err := t.Resolve(pkg)
	if err != nil {
		fmt.Printf("'%v': FATAL: %v\n", pkg, err)
		return
	}
	if outputJSON {
		writePkgJSON(os.Stdout, *t.Root)
	} else {
		writePkg(os.Stdout, *t.Root, 0, false)
	}
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
