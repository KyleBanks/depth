package benchmark

import (
	"testing"

	"github.com/KyleBanks/depth"
)

func BenchmarkTree_ResolveStrings(b *testing.B) {
	var t depth.Tree
	benchmarkTreeResolveStrings(&t, b)
}

func BenchmarkTree_ResolveStringsInternal(b *testing.B) {
	benchmarkTreeResolveStrings(&depth.Tree{
		ResolveInternal: true,
	}, b)
}

func BenchmarkTree_ResolveStringsTest(b *testing.B) {
	benchmarkTreeResolveStrings(&depth.Tree{
		ResolveTest: true,
	}, b)
}

func BenchmarkTree_ResolveStringsInternalTest(b *testing.B) {
	benchmarkTreeResolveStrings(&depth.Tree{
		ResolveInternal: true,
		ResolveTest:     true,
	}, b)
}

func benchmarkTreeResolveStrings(t *depth.Tree, b *testing.B) {
	for i := 0; i < b.N; i++ {
		if err := t.Resolve("strings"); err != nil {
			b.Fatal(err)
		}
	}
}
