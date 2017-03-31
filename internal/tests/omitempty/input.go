// Copyright 2017 Felix Lange <fjl@twurst.com>.
// Use of this source code is governed by the MIT license,
// which can be found in the LICENSE file.

//go:generate gencodec -type X -field-override Xo -formats json,yaml,toml -out output.go

package omitempty

type replacedInt int

type X struct {
	Int int `json:",omitempty"`
}

type Xo struct {
	Int replacedInt
}
