// Copyright 2017 Felix Lange <fjl@twurst.com>.
// Use of this source code is governed by the MIT license,
// which can be found in the LICENSE file.

//go:generate gencodec -type X -field-override Xo -formats json,yaml -out output.go

package mapconv

type replacedString string

type replacedInt int

type X struct {
	M map[string]int
}

type Xo struct {
	M map[replacedString]replacedInt
}
