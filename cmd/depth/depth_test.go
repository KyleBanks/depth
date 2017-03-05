package main

import (
	"github.com/KyleBanks/depth"
)

func Example_handlePkgsStrings() {
	var t depth.Tree

	handlePkgs(&t, []string{"strings"})
	// Output:
	// strings
	//   ├ errors
	//   ├ io
	//   ├ unicode
	//   └ unicode/utf8
}

func Example_handlePkgsTestStrings() {
	var t depth.Tree
	t.ResolveTest = true

	handlePkgs(&t, []string{"strings"})
	// Output:
	// strings
	//   ├ bytes
	//   ├ errors
	//   ├ fmt
	//   ├ io
	//   ├ io/ioutil
	//   ├ math/rand
	//   ├ reflect
	//   ├ sync
	//   ├ testing
	//   ├ unicode
	//   ├ unicode/utf8
	//   └ unsafe
}

func Example_handlePkgsUnknown() {
	var t depth.Tree

	handlePkgs(&t, []string{"notreal"})
	// Output:
	// 'notreal': FATAL: unable to resolve root package
}
