// Copyright 2017 Felix Lange <fjl@twurst.com>.
// Use of this source code is governed by the MIT license,
// which can be found in the LICENSE file.

//go:generate gencodec -type Y -field-override Yo -formats json,yaml,toml -out output.go

package nameclash

import (
	errors "github.com/fjl/gencodec/internal/clasherrors"
	json "github.com/fjl/gencodec/internal/clashjson"
)

// This one clashes with the generated intermediate type name.
type YJSON struct{}

// This type clashes with a name in the override struct.
type enc int

// These types clash with variables, but are ignored because they're
// not referenced by the struct.
type input struct{}
type dec struct{}

type Y struct {
	Foo    json.Foo   `optional:"true"`
	Foo2   json.Foo   `optional:"true"`
	Bar    errors.Foo `optional:"true"`
	Gazonk YJSON      `optional:"true"`
	Over   int        `optional:"true"`
}

type Yo struct {
	Over enc
}
