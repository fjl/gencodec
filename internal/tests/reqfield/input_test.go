// Copyright 2017 Felix Lange <fjl@twurst.com>.
// Use of this source code is governed by the MIT license,
// which can be found in the LICENSE file.

package reqfield

import (
	"encoding/json"
	"testing"
)

func TestRequiredJSON(t *testing.T) {
	input := `{"required": 1}`
	var x X
	if err := json.Unmarshal([]byte(input), &x); err != nil {
		t.Fatal("unexpected error", err)
	}

	input = `{}`
	if err := json.Unmarshal([]byte(input), &x); err == nil {
		t.Fatal("expected error, got nil")
	}
}
