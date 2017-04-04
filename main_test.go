// Copyright 2017 Felix Lange <fjl@twurst.com>.
// Use of this source code is governed by the MIT license,
// which can be found in the LICENSE file.

package main

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/kylelemons/godebug/diff"
)

// 'golden' tests. These tests compare the output code with the expected
// code in internal/tests/*/output.go. The expected can be updated using
//
//    go generate ./internal/...

func TestGolden(t *testing.T) {
	tests := []Config{
		Config{Dir: "mapconv", Type: "X", FieldOverride: "Xo", Formats: AllFormats},
		Config{Dir: "sliceconv", Type: "X", FieldOverride: "Xo", Formats: AllFormats},
		Config{Dir: "nameclash", Type: "Y", FieldOverride: "Yo", Formats: AllFormats},
		Config{Dir: "omitempty", Type: "X", FieldOverride: "Xo", Formats: AllFormats},
		Config{Dir: "reqfield", Type: "X", Formats: []string{"json"}},
	}
	for _, test := range tests {
		test := test
		t.Run(test.Dir, func(t *testing.T) {
			t.Parallel()
			runGoldenTest(t, test)
		})
	}
}

func runGoldenTest(t *testing.T, cfg Config) {
	cfg.Dir = filepath.Join("internal", "tests", cfg.Dir)
	want, err := ioutil.ReadFile(filepath.Join(cfg.Dir, "output.go"))
	if err != nil {
		t.Fatal(err)
	}

	code, err := cfg.process()
	if err != nil {
		t.Fatal(err)
	}
	if d := diff.Diff(string(want), string(code)); d != "" {
		t.Errorf("output mismatch\n\n%s", d)
	}
}
