package depth

import (
	"testing"
)

func BenchmarkTree_ResolveStrings(b *testing.B) {
	var t Tree
	benchmarkTreeResolveStrings(t, b)
}

func BenchmarkTree_ResolveStringsInternal(b *testing.B) {
	var t Tree
	t.ResolveInternal = true
	benchmarkTreeResolveStrings(t, b)
}

func BenchmarkTree_ResolveStringsTest(b *testing.B) {
	var t Tree
	t.ResolveTest = true
	benchmarkTreeResolveStrings(t, b)
}

func BenchmarkTree_ResolveStringsInternalTest(b *testing.B) {
	var t Tree
	t.ResolveInternal = true
	t.ResolveTest = true
	benchmarkTreeResolveStrings(t, b)
}

func benchmarkTreeResolveStrings(t Tree, b *testing.B) {
	for i := 0; i < b.N; i++ {
		if err := t.Resolve("strings"); err != nil {
			b.Fatal(err)
		}
	}
}
