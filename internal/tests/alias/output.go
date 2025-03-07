// Code generated by github.com/fjl/gencodec. DO NOT EDIT.

package alias

import (
	"encoding/json"

	"github.com/fjl/gencodec/internal/tests/alias/other"
)

// MarshalJSON marshals as JSON.
func (x X) MarshalJSON() ([]byte, error) {
	type X struct {
		A Aliased
		B AliasedTwice
		C other.Int
	}
	var enc X
	enc.A = x.A
	enc.B = x.B
	enc.C = x.C
	return json.Marshal(&enc)
}

// UnmarshalJSON unmarshals from JSON.
func (x *X) UnmarshalJSON(input []byte) error {
	type X struct {
		A *Aliased
		B *AliasedTwice
		C *other.Int
	}
	var dec X
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	if dec.A != nil {
		x.A = *dec.A
	}
	if dec.B != nil {
		x.B = *dec.B
	}
	if dec.C != nil {
		x.C = *dec.C
	}
	return nil
}
