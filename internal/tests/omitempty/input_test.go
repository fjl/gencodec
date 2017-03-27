// Copyright 2017 Felix Lange <fjl@twurst.com>.
// Use of this source code is governed by the MIT license,
// which can be found in the LICENSE file.

package omitempty

import (
	"encoding/json"
	"testing"
)

func TestOmitemptyJSON(t *testing.T) {
	want := `{}`
	out, err := json.Marshal(new(X))
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != want {
		t.Fatalf("got %#q, want %#q", string(out), want)
	}
}
