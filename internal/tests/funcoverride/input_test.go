// Copyright 2017 Felix Lange <fjl@twurst.com>.
// Use of this source code is governed by the MIT license,
// which can be found in the LICENSE file.

package funcoverride

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestOverrideFuncJSON(t *testing.T) {
	z := Z{"str", 1234}
	hash := z.Hash()
	multiply := z.MultiplyIByTwo()
	want := fmt.Sprintf(`{"s":"%s","iVal":%d,"Hash":"%s","multipliedByTwo":%d}`, z.S, z.I, hash, multiply)
	out, err := json.Marshal(z)
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != want {
		t.Fatalf("got %#q, want %#q", string(out), want)
	}

	var zUnmarshaled Z
	if err := json.Unmarshal([]byte(want), &zUnmarshaled); err != nil {
		t.Fatalf("could not unmarshal Z: %v", err)
	}
	if zUnmarshaled.I != z.I {
		t.Errorf("Z.I has an unexpected value, want %d, got %d", z.I, zUnmarshaled.I)
	}
	if zUnmarshaled.S != z.S {
		t.Errorf("Z.Str has an unexpected value, want %s, got %s", z.S, zUnmarshaled.S)
	}
	uHash := zUnmarshaled.Hash()
	if uHash != hash {
		t.Errorf("Z.Hash() returned unexpected value, want %s, got %s", hash, uHash)
	}
	uMultiply := zUnmarshaled.MultiplyIByTwo()
	if uMultiply != multiply {
		t.Errorf("Z.MultiplIByTwo() returned unexpected value, want %d, got %d", multiply, uMultiply)
	}
}
