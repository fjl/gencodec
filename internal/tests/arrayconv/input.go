// Copyright 2017 Felix Lange <fjl@twurst.com>.
// Use of this source code is governed by the MIT license,
// which can be found in the LICENSE file.

//go:generate go run github.com/fjl/gencodec -type X -field-override Xo -formats json -out output.go

package arrayconv

type MyArray [32]int

type X struct {
	A MyArray
	RequiredA MyArray `gencodec:"required"`
}

type Xo struct{
	A []int
	RequiredA []int
}
