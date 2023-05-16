// Copyright 2017 Felix Lange <fjl@twurst.com>.
// Use of this source code is governed by the MIT license,
// which can be found in the LICENSE file.

//go:generate go run github.com/fjl/gencodec -type X -field-override Xo -formats json,yaml,toml -out output.go

package mapconv

type replacedString string

type replacedInt int

type namedMap map[string]int

type namedMap2 map[replacedString]replacedInt

type X struct {
	Map         map[string]int
	Named       namedMap
	NoConv      map[string]int
	NoConvNamed map[string]int
}

func (x *X) Func() map[string]int {
	return map[string]int{"a": 1, "b": 2}
}

type Xo struct {
	Map         map[replacedString]replacedInt
	Named       namedMap2
	NoConvNamed namedMap

	Func map[replacedString]replacedInt
}
