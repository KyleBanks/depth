package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"log"
	"bytes"
	"strings"
	"github.com/goccy/go-graphviz"

	"github.com/KyleBanks/depth"
)

const (
	outputClosedPadding = "  "
	outputOpenPadding   = "│ "
	outputPrefix        = "├ "
	outputPrefixLast    = "└ "
)

var outputJSON bool
var explainPkg string
var outputGraph bool

type summary struct {
	numInternal int
	numExternal int
	numTesting  int
}

func main() {
	t, pkgs := parse(os.Args[1:])
	if err := handlePkgs(t, pkgs, outputJSON, explainPkg,outputGraph); err != nil {
		os.Exit(1)
	}
}

// parse constructs a depth.Tree from command-line arguments, and returns the
// remaining user-supplied package names
func parse(args []string) (*depth.Tree, []string) {
	f := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	var t depth.Tree
	f.BoolVar(&t.ResolveInternal, "internal", false, "If set, resolves dependencies of internal (stdlib) packages.")
	f.BoolVar(&t.ResolveTest, "test", false, "If set, resolves dependencies used for testing.")
	f.IntVar(&t.MaxDepth, "max", 0, "Sets the maximum depth of dependencies to resolve.")
	f.BoolVar(&outputJSON, "json", false, "If set, outputs the depencies in JSON format.")
	f.StringVar(&explainPkg, "explain", "", "If set, show which packages import the specified target")
	f.BoolVar(&outputGraph, "graph", false, "If set, renders the dependencies in a graph and outputs as a png file")
	f.Parse(args)

	return &t, f.Args()
}

// handlePkgs takes a slice of package names, resolves a Tree on them,
// and outputs each Tree to Stdout.
func handlePkgs(t *depth.Tree, pkgs []string, outputJSON bool, explainPkg string, outputGraph bool) error {
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

		if explainPkg != "" {
			writeExplain(os.Stdout, *t.Root, []string{}, explainPkg)
			continue
		}

		if outputGraph {
			generateGraph(*t.Root)
			continue
		}

		writePkg(os.Stdout, *t.Root)
		writePkgSummary(os.Stdout, *t.Root)
	}
	return nil
}

// writePkgSummary writes a summary of all packages in a tree
func writePkgSummary(w io.Writer, pkg depth.Pkg) {
	var sum summary
	set := make(map[string]struct{})
	for _, p := range pkg.Deps {
		collectSummary(&sum, p, set)
	}
	fmt.Fprintf(w, "%d dependencies (%d internal, %d external, %d testing).\n",
		sum.numInternal+sum.numExternal,
		sum.numInternal,
		sum.numExternal,
		sum.numTesting)
}

func collectSummary(sum *summary, pkg depth.Pkg, nameSet map[string]struct{}) {
	if _, ok := nameSet[pkg.Name]; !ok {
		nameSet[pkg.Name] = struct{}{}
		if pkg.Internal {
			sum.numInternal++
		} else {
			sum.numExternal++
		}
		if pkg.Test {
			sum.numTesting++
		}
		for _, p := range pkg.Deps {
			collectSummary(sum, p, nameSet)
		}
	}
}

// writePkgJSON writes the full Pkg as JSON to the provided Writer.
func writePkgJSON(w io.Writer, p depth.Pkg) {
	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	e.Encode(p)
}

func writePkg(w io.Writer, p depth.Pkg) {
	fmt.Fprintf(w, "%s\n", p.String())

	for idx, d := range p.Deps {
		writePkgRec(w, d, []bool{true}, idx == len(p.Deps)-1)
	}
}

// writePkg recursively prints a Pkg and its dependencies to the Writer provided.
func writePkgRec(w io.Writer, p depth.Pkg, closed []bool, isLast bool) {
	var prefix string

	for _, c := range closed {
		if c {
			prefix += outputClosedPadding
			continue
		}

		prefix += outputOpenPadding
	}

	closed = append(closed, false)
	if isLast {
		prefix += outputPrefixLast
		closed[len(closed)-1] = true
	} else {
		prefix += outputPrefix
	}

	fmt.Fprintf(w, "%v%v\n", prefix, p.String())

	for idx, d := range p.Deps {
		writePkgRec(w, d, closed, idx == len(p.Deps)-1)
	}
}

// writeExplain shows possible paths for a given package.
func writeExplain(w io.Writer, pkg depth.Pkg, stack []string, explain string) {
	stack = append(stack, pkg.Name)
	if pkg.Name == explain {
		fmt.Fprintln(w, strings.Join(stack, " -> "))
	}
	for _, p := range pkg.Deps {
		writeExplain(w, p, stack, explain)
	}
}

// generateGraph creates a node graph and renders it as a file
func generateGraph(pkg depth.Pkg) {
	
	//adjacency list
adjlist := make(map[string][]string)
for _, p := range pkg.Deps {

	_, parentok := adjlist[p.Parent.Name]
	if !parentok {
		adjlist[p.Parent.Name] = make([]string, 0)
	}
	adjlist[p.Parent.Name] = append(adjlist[p.Parent.Name], p.Name)

	_, ok := adjlist[p.Name]
	if !ok {
		adjlist[p.Name] = make([]string, 0)
	}


	generateGraphRec(p, p.Name, adjlist)
}

g := graphviz.New()
graph, err := g.Graph()

defer func() {
	if err := graph.Close(); err != nil {
		log.Fatal(err)
	}
	g.Close()
}()

if err != nil {
	log.Fatal(err)
}

for parent, children := range adjlist {
	n, err := graph.CreateNode(parent)
	if err != nil {
		log.Fatal(err)
	}

	for _, child := range children {
		c, err := graph.CreateNode(child)
		if err != nil {
			log.Fatal(err)
		}

		_, err = graph.CreateEdge("", n, c)
		if err != nil {
			log.Fatal(err)
		}

	}
}

var buf bytes.Buffer
if err := g.Render(graph, "dot", &buf); err != nil {
	log.Fatal(err)
}

if err := g.RenderFilename(graph, graphviz.PNG, "graph.png"); err != nil {
	log.Fatal(err)
}
}

func generateGraphRec(pkg depth.Pkg, parent string, adjlist map[string][]string) {

for _, p := range pkg.Deps {
	_, ok := adjlist[p.Name]
	if !ok {
		adjlist[p.Name] = make([]string, 0)
	}

	adjlist[parent] = append(adjlist[parent], p.Name)
	generateGraphRec(p, p.Name, adjlist)
}
}