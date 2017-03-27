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

func TestMapConv(t *testing.T) {
	runGoldenTest(t, Config{Dir: "mapconv", Type: "X", FieldOverride: "Xo", Formats: AllFormats})
}

func TestSliceConv(t *testing.T) {
	runGoldenTest(t, Config{Dir: "sliceconv", Type: "X", FieldOverride: "Xo", Formats: AllFormats})
}

func TestNameClash(t *testing.T) {
	runGoldenTest(t, Config{Dir: "nameclash", Type: "Y", FieldOverride: "Yo", Formats: AllFormats})
}

func TestOmitempty(t *testing.T) {
	runGoldenTest(t, Config{Dir: "omitempty", Type: "X", FieldOverride: "Xo", Formats: AllFormats})
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
