// Copyright 2017 Felix Lange <fjl@twurst.com>.
// Use of this source code is governed by the MIT license,
// which can be found in the LICENSE file.

//go:generate gencodec -type X -field-override Xo -formats json,yaml -out output.go

package sliceconv

type replacedInt int

type namedSlice []int

type namedSlice2 []replacedInt

type X struct {
	Slice []int
	Named namedSlice
}

type Xo struct {
	Slice []replacedInt
	Named namedSlice2
}
