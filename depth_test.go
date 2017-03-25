package depth

import (
	"go/build"
	"testing"
)

type MockImporter struct {
	ImportFn func(name, srcDir string, im build.ImportMode) (*build.Package, error)
}

func (m MockImporter) Import(name, srcDir string, im build.ImportMode) (*build.Package, error) {
	return m.ImportFn(name, srcDir, im)
}

func TestTree_Resolve(t *testing.T) {
	// Fail case, bad package name
	var tr Tree
	if err := tr.Resolve("name"); err != ErrRootPkgNotResolved {
		t.Fatalf("Unexpected error, expected=%v, got=%b", ErrRootPkgNotResolved, err)
	}

	// Positive case, expect deps
	if err := tr.Resolve("strings"); err != nil {
		t.Fatal(err)
	}

	if tr.Root == nil || tr.Root.Name != "strings" {
		t.Fatalf("Unexpected Root, expected=%v, got=%v", "strings", tr.Root)
	} else if len(tr.Root.Deps) == 0 {
		t.Fatal("Expected positive number of Deps")
	} else if len(tr.Root.SrcDir) == 0 {
		t.Fatal("Expected SrcDir to be populated")
	} else if tr.Root.Raw == nil {
		t.Fatal("Expected non-nil Raw")
	}

	// Reuse the same tree and the same package to ensure that the internal pkg cache
	// is reset and dependencies are still resolved.
	stringsDepCount := len(tr.Root.Deps)
	if err := tr.Resolve("strings"); err != nil {
		t.Fatal(err)
	}
	if len(tr.Root.Deps) != stringsDepCount {
		t.Fatalf("Unexpected number of Deps, expected=%v, got=%b", stringsDepCount, len(tr.Root.Deps))
	}
}

func TestTree_shouldResolveInternal(t *testing.T) {
	var pt Tree
	pt.Root = &Pkg{}

	if pt.shouldResolveInternal(&Pkg{}) {
		t.Fatal("Unexpected shouldResolveInternal, should have been false for non-root pkg and default config")
	}

	pt.ResolveInternal = true
	if !pt.shouldResolveInternal(&Pkg{}) {
		t.Fatal("Unexpected shouldResolveInternal, should have been true when ResolveInternal = true")
	}
	pt.ResolveInternal = false

	if !pt.shouldResolveInternal(pt.Root) {
		t.Fatal("Unexpected shouldResolveInternal, should have been true for root pkg")
	}
}

func TestTree_isAtMaxDepth(t *testing.T) {
	tests := []struct {
		maxDepth int
		depth    int
		expected bool
	}{
		{0, 0, false},
		{0, 10, false},
		{1, 0, false},
		{1, 1, true},
		{1, 10, true},
	}

	for idx, tt := range tests {
		tr := Tree{MaxDepth: tt.maxDepth}

		var last *Pkg
		for i := 0; i < tt.depth+1; i++ {
			p := Pkg{Parent: last}
			last = &p
		}

		maxDepth := tr.isAtMaxDepth(last)
		if maxDepth != tt.expected {
			t.Fatalf("[%v] Unexpected isAtMaxDepth, expected=%v, got=%v", idx, tt.expected, maxDepth)
		}
	}
}

func TestTree_hasSeenImport(t *testing.T) {
	var tr Tree

	if tr.hasSeenImport("name") {
		t.Fatalf("Expected false the first time an import name is provided, got=true")
	}

	if !tr.hasSeenImport("name") {
		t.Fatalf("Expected true to be returned after the import name has been seen, got=false")
	}
}
