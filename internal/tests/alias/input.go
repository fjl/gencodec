// Copyright 2025 Felix Lange <fjl@twurst.com>.
// Use of this source code is governed by the MIT license,
// which can be found in the LICENSE file.

//go:generate go run github.com/fjl/gencodec -type X -out output.go

package alias

import (
	"math/big"
)

// Alias types chosen because their originals have special handling that is easy
// to spot when inspecting generated output.
type (
	Aliased = big.Int
)

type X struct {
	A Aliased
}
