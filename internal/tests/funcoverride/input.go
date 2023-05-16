// Copyright 2017 Felix Lange <fjl@twurst.com>.
// Use of this source code is governed by the MIT license,
// which can be found in the LICENSE file.

//go:generate go run github.com/fjl/gencodec -type Z -field-override Zo -formats json,yaml,toml -out output.go

package funcoverride

import (
	"fmt"
)

type Z struct {
	S string `json:"s"`
	I int32  `json:"iVal"`
}

func (z *Z) Hash() string {
	return fmt.Sprintf("%s-%d", z.S, z.I)
}

func (z *Z) MultiplyIByTwo() int32 {
	return 2 * z.I
}

func (z *Z) NotUsed() string {
	return "not used"
}

type Zo struct {
	Hash           string
	MultiplyIByTwo int64 `json:"multipliedByTwo"`
}
