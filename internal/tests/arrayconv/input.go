// Copyright 2017 Felix Lange <fjl@twurst.com>.
// Use of this source code is governed by the MIT license,
// which can be found in the LICENSE file.

//go:generate go run github.com/fjl/gencodec -type X -field-override Xo -formats json -out output.go

package arrayconv

type MyArray [32]int

type MyInt int

type MySlice []int64

type X struct {
	A         MyArray
	A2        MyArray
	RequiredA MyArray `gencodec:"required"`

	S         MySlice
	RequiredS MySlice `gencodec:"required"`
}

type Xo struct {
	A         []int
	RequiredA []int
	A2        []MyInt
	S         [16]int64
	RequiredS [16]int64
}
