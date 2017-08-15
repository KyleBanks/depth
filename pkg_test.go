package depth

import (
	"go/build"
	"sort"
	"testing"
)

func TestPkg_CleanName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"strings", "strings"},
		{"net/http", "net/http"},
		{"github.com/KyleBanks/depth", "github.com/KyleBanks/depth"},
		{"C", ""},
		{"golang_org/x/anything", "vendor/golang_org/x/anything"},
	}

	for _, tt := range tests {
		p := Pkg{Name: tt.input}

		out := p.cleanName()
		if out != tt.expected {
			t.Fatalf("Unexpected cleanName, expected=%v, got=%v", tt.expected, out)
		}
	}
}

func TestPkg_AddDepImportSeen(t *testing.T) {
	var m MockImporter
	var tr Tree
	tr.Importer = m

	testName := "test"
	testSrcDir := "src/testing"
	var expectedIm build.ImportMode

	p := Pkg{Tree: &tr}
	m.ImportFn = func(name, srcDir string, im build.ImportMode) (*build.Package, error) {
		if name != testName {
			t.Fatalf("Unexpected name provided, expected=%v, got=%v", testName, name)
		}
		if srcDir != testSrcDir {
			t.Fatalf("Unexpected srcDir provided, expected=%v, got=%v", testSrcDir, srcDir)
		}
		if im != expectedIm {
			t.Fatalf("Unexpected ImportMode provided, expected=%v, got=%v", expectedIm, im)
		}

		return &build.Package{}, nil
	}

	// Hasn't seen the import
	p.addDep(m, testName, testSrcDir, false)

	// Has seen the import
	expectedIm = build.FindOnly
	p.addDep(m, testName, testSrcDir, false)
}

func TestByInternalAndName(t *testing.T) {
	pkgs := []Pkg{
		Pkg{Internal: true, Name: "net/http"},
		Pkg{Internal: false, Name: "github.com/KyleBanks/depth"},
		Pkg{Internal: true, Name: "strings"},
		Pkg{Internal: false, Name: "github.com/KyleBanks/commuter"},
	}
	expected := []string{"net/http", "strings", "github.com/KyleBanks/commuter", "github.com/KyleBanks/depth"}

	sort.Sort(byInternalAndName(pkgs))

	for i, e := range expected {
		if pkgs[i].Name != e {
			t.Fatalf("Unexpected Pkg at index %v, expected=%v, got=%v", i, e, pkgs[i].Name)
		}
	}
}
