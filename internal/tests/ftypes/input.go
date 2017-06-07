// Copyright 2017 Felix Lange <fjl@twurst.com>.
// Use of this source code is governed by the MIT license,
// which can be found in the LICENSE file.

//go:generate gencodec -type X -formats json -out output.go

package ftypes

type X struct {
	Int int
	Err error
	If  interface{}
}
