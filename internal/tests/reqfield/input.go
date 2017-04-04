// Copyright 2017 Felix Lange <fjl@twurst.com>.
// Use of this source code is governed by the MIT license,
// which can be found in the LICENSE file.

//go:generate gencodec -type X -formats json -out output.go

package reqfield

type X struct {
	Required    int `gencodec:"required"`
	NotRequired int `gencodec:"required" json:"-"`
}
