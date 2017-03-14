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

func Example_handlePkgsDepth() {
	var t depth.Tree

	handlePkgs(&t, []string{"github.com/KyleBanks/depth/cmd/depth"})
	// Output:
	// github.com/KyleBanks/depth/cmd/depth
	//   ├ encoding/json
	//   ├ flag
	//   ├ fmt
	//   ├ io
	//   ├ os
	//   ├ strings
	//   └ github.com/KyleBanks/depth
	//     ├ bytes
	//     ├ errors
	//     ├ go/build
	//     ├ os
	//     ├ path
	//     ├ sort
	//     └ strings
}

func Example_handlePkgsUnknown() {
	var t depth.Tree

	handlePkgs(&t, []string{"notreal"})
	// Output:
	// 'notreal': FATAL: unable to resolve root package
}

func Example_handlePkgsJson() {
	outputJSON = true
	var t depth.Tree
	handlePkgs(&t, []string{"strings"})

	// Output:
	// {
	//   "name": "strings",
	//   "resolved": true,
	//   "deps": [
	//     {
	//       "name": "errors",
	//       "resolved": true,
	//       "deps": null
	//     },
	//     {
	//       "name": "io",
	//       "resolved": true,
	//       "deps": null
	//     },
	//     {
	//       "name": "unicode",
	//       "resolved": true,
	//       "deps": null
	//     },
	//     {
	//       "name": "unicode/utf8",
	//       "resolved": true,
	//       "deps": null
	//     }
	//   ]
	// }
}
