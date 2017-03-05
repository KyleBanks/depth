package depth

import (
	"go/build"
	"path"
	"sort"
	"strings"
	"sync"
)

// Pkg represents a Go source package, and its dependencies.
type Pkg struct {
	Name     string `json:"name"`
	SrcDir   string `json:"-"`
	Internal bool   `json:"-"`

	Tree   *Tree `json:"-"`
	Parent *Pkg  `json:"-"`
	mu     *sync.Mutex
	Deps   []Pkg `json:"deps"`
}

// Resolve recursively finds all dependencies for the Pkg and the packages it depends on.
func (p *Pkg) Resolve(i Importer, resolveImports bool) error {
	name := p.cleanName()
	if name == "" {
		return nil
	}

	// Stop resolving imports if we've reached max depth.
	if resolveImports && p.Tree.isAtMaxDepth(p) {
		resolveImports = false
	}

	var importMode build.ImportMode
	if !resolveImports {
		importMode = build.FindOnly
	}

	pkg, err := i.Import(name, p.SrcDir, importMode)
	if err != nil {
		return err
	}

	// If this is an internal dependency, we may need to skip it.
	if pkg.Goroot {
		p.Internal = true

		if !p.Tree.shouldResolveInternal(p) {
			return nil
		}
	}

	imports := pkg.Imports
	if p.Tree.ResolveTest {
		imports = append(imports, append(pkg.TestImports, pkg.XTestImports...)...)
	}

	return p.setDeps(i, imports, pkg.Dir)
}

// setDeps takes a slice of import paths and the source directory they are relative to,
// and creates the Deps of the Pkg. Each dependency is also further resolved prior to being added
// to the Pkg.
func (p *Pkg) setDeps(i Importer, imports []string, srcDir string) error {
	p.mu = &sync.Mutex{}
	unique := make(map[string]struct{})
	errCh := make(chan error)
	var wg sync.WaitGroup

	for _, imp := range imports {
		// Mostly for testing files where cyclic imports are allowed.
		if imp == p.Name {
			continue
		}
		// Skip duplicates.
		if _, ok := unique[imp]; ok {
			continue
		}
		unique[imp] = struct{}{}

		wg.Add(1)
		go func(imp string) {
			err := p.addDep(i, imp, srcDir)

			wg.Done()
			if err != nil {
				errCh <- err
			}
		}(imp)
	}

	wg.Wait()

	select {
	case err := <-errCh:
		return err
	default:
		sort.Sort(byInternalAndName(p.Deps))
		return nil
	}
}

// addDep creates a Pkg and it's dependencies from an imported package name.
func (p *Pkg) addDep(i Importer, name string, srcDir string) error {
	// Don't resolve imports for Pkgs that we've already seen and resolved.
	resolveImports := !p.Tree.hasSeenImport(name)

	dep := Pkg{
		Name:   name,
		SrcDir: srcDir,
		Tree:   p.Tree,
		Parent: p,
	}

	if err := dep.Resolve(i, resolveImports); err != nil {
		return err
	}

	p.mu.Lock()
	p.Deps = append(p.Deps, dep)
	p.mu.Unlock()

	return nil
}

// isParent goes recursively up the chain of Pkgs to determine if the name provided is ever a
// parent of the current Pkg.
func (p *Pkg) isParent(name string) bool {
	if p.Parent == nil {
		return false
	}

	if p.Parent.Name == name {
		return true
	}

	return p.Parent.isParent(name)
}

// depth returns the depth of the Pkg within the Tree.
func (p *Pkg) depth() int {
	if p.Parent == nil {
		return 0
	}

	return p.Parent.depth() + 1
}

// cleanName returns a cleaned version of the Pkg name used for resolving dependencies.
//
// If an empty string is returned, dependencies should not be resolved.
func (p *Pkg) cleanName() string {
	name := p.Name

	// C 'package' cannot be resolved.
	if name == "C" {
		return ""
	}

	// Internal golang_org/* packages must be prefixed with vendor/
	//
	// Thanks to @davecheney for this:
	// https://github.com/davecheney/graphpkg/blob/master/main.go#L46
	if strings.HasPrefix(name, "golang_org") {
		name = path.Join("vendor", name)
	}

	return name
}

// byInternalAndName ensures a slice of Pkgs are sorted such that the internal stdlib
// packages are always above external packages (ie. github.com/whatever).
type byInternalAndName []Pkg

func (b byInternalAndName) Len() int {
	return len(b)
}

func (b byInternalAndName) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b byInternalAndName) Less(i, j int) bool {
	if b[i].Internal && !b[j].Internal {
		return true
	} else if !b[i].Internal && b[j].Internal {
		return false
	}

	return b[i].Name < b[j].Name
}
