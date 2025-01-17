// Copyright 2025 Felix Lange <fjl@twurst.com>.
// Use of this source code is governed by the MIT license,
// which can be found in the LICENSE file.

//go:generate go run github.com/fjl/gencodec -type Cfg -field-override cfgOverride -formats json,yaml,toml -out output.go

package ifaceoverride

type Iface interface {
	Method() string
}

type Impl struct {
}

func (*Impl) Method() string {
	return "yes"
}

type Cfg struct {
	Field Iface
}

type cfgOverride struct {
	Field *Impl
}
