// Copyright 2017 Felix Lange <fjl@twurst.com>.
// Use of this source code is governed by the MIT license,
// which can be found in the LICENSE file.

//go:generate gencodec -type X -field-override Xo -formats json,yaml -out output.go

package mapconv

type replacedString string

type replacedInt int

type namedMap map[string]int

type namedMap2 map[replacedString]replacedInt

type X struct {
	Map   map[string]int
	Named namedMap
}

type Xo struct {
	Map   map[replacedString]replacedInt
	Named namedMap2
}
