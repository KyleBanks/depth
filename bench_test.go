package depth

import (
	"testing"
)

func BenchmarkTree_ResolveStrings(b *testing.B) {
	benchmarkTreeResolveStrings(&Tree{}, b)
}

func BenchmarkTree_ResolveStringsInternal(b *testing.B) {
	benchmarkTreeResolveStrings(&Tree{
		ResolveInternal: true,
	}, b)
}

func BenchmarkTree_ResolveStringsTest(b *testing.B) {
	benchmarkTreeResolveStrings(&Tree{
		ResolveTest: true,
	}, b)
}

func BenchmarkTree_ResolveStringsInternalTest(b *testing.B) {
	benchmarkTreeResolveStrings(&Tree{
		ResolveInternal: true,
		ResolveTest:     true,
	}, b)
}

func benchmarkTreeResolveStrings(t *Tree, b *testing.B) {
	for i := 0; i < b.N; i++ {
		if err := t.Resolve("strings"); err != nil {
			b.Fatal(err)
		}
	}
}
