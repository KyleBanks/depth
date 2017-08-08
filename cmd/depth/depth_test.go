package main

import (
	"fmt"
	"testing"

	"github.com/KyleBanks/depth"
)

func Test_parse(t *testing.T) {
	tests := []struct {
		internal bool
		test     bool
		depth    int
		json     bool
	}{
		{true, true, 0, true},
		{false, false, 10, false},
		{true, false, 10, false},
		{false, true, 5, true},
	}

	for idx, tt := range tests {
		tr, _ := parse([]string{
			fmt.Sprintf("-internal=%v", tt.internal),
			fmt.Sprintf("-test=%v", tt.test),
			fmt.Sprintf("-max=%v", tt.depth),
			fmt.Sprintf("-json=%v", tt.json),
		})

		if tr.ResolveInternal != tt.internal {
			t.Fatalf("[%v] Unexpected ResolveInternal, expected=%v, got=%v", idx, tt.internal, tr.ResolveInternal)
		} else if tr.ResolveTest != tt.test {
			t.Fatalf("[%v] Unexpected ResolveTest, expected=%v, got=%v", idx, tt.test, tr.ResolveTest)
		} else if tr.MaxDepth != tt.depth {
			t.Fatalf("[%v] Unexpected MaxDepth, expected=%v, got=%v", idx, tt.depth, tr.MaxDepth)
		} else if outputJSON != tt.json {
			t.Fatalf("[%v] Unexpected outputJSON, expected=%v, got=%v", idx, tt.json, outputJSON)
		}
	}
}

func Example_handlePkgsStrings() {
	var t depth.Tree

	handlePkgs(&t, []string{"strings"}, false)
	// Output:
	// strings
	//   ├ errors
	//   ├ io
	//   ├ unicode
	//   └ unicode/utf8
	// 4 dependencies (4 internal, 0 external, 0 testing).

}

func Example_handlePkgsTestStrings() {
	var t depth.Tree
	t.ResolveTest = true

	handlePkgs(&t, []string{"strings"}, false)
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
	// 12 dependencies (12 internal, 0 external, 8 testing).
}

func Example_handlePkgsDepth() {
	var t depth.Tree

	handlePkgs(&t, []string{"github.com/KyleBanks/depth/cmd/depth"}, false)
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
	// 12 dependencies (11 internal, 1 external, 0 testing).
}

func Example_handlePkgsUnknown() {
	var t depth.Tree

	handlePkgs(&t, []string{"notreal"}, false)
	// Output:
	// 'notreal': FATAL: unable to resolve root package
}

func Example_handlePkgsJson() {
	var t depth.Tree
	handlePkgs(&t, []string{"strings"}, true)

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
