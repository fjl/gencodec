// Copyright 2025 Felix Lange <fjl@twurst.com>.
// Use of this source code is governed by the MIT license,
// which can be found in the LICENSE file.

//go:generate go run github.com/fjl/gencodec -type X -field-override xOverride -out output.go

package alias

import (
	"math/big"

	"github.com/fjl/gencodec/internal/tests/alias/other"
)

// Alias types chosen because their originals have special handling that is easy
// to spot when inspecting generated output.
type (
	Aliased = big.Int
	// Demonstrate recursive unaliasing
	intermediate = big.Int
	AliasedTwice = intermediate

	Element      struct{}
	ElementDeriv Element
	Slice        []Element
	SliceAlias   = Slice
)

type X struct {
	A Aliased
	B AliasedTwice
	C other.Int
	D SliceAlias
}

type xOverride struct {
	D []ElementDeriv
}
