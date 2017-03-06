// Copyright 2017 Felix Lange <fjl@twurst.com>.
// Use of this source code is governed by the MIT license,
// which can be found in the LICENSE file.

//go:generate gencodec -type X -field-override Xo -out output.go

package sliceconv

type replacedInt int

type X struct {
	S []int
}

type Xo struct {
	S []replacedInt
}
