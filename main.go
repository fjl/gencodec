// Copyright 2017 Felix Lange <fjl@twurst.com>.
// Use of this source code is governed by the MIT license,
// which can be found in the LICENSE file.

/*
Command gencodec generates marshaling methods for struct types.

When gencodec is invoked on a directory and type name, it creates a Go source file
containing JSON and YAML marshaling methods for the type. The generated methods add
features which the standard json package cannot offer.

	gencodec -dir . -type MyType -out mytype_json.go

Struct Tags

All fields are required unless the "optional" struct tag is present. The generated
unmarshaling method returns an error if a required field is missing. Other struct tags are
carried over as is. The standard "json" and "yaml" tags can be used to rename a field when
marshaling.

Example:

	type foo struct {
		Required string
		Optional string `optional:""`
		Renamed  string `json:"otherName"`
	}

Field Type Overrides

An invocation of gencodec can specify an additional 'field override' struct from which
marshaling type replacements are taken. If the override struct contains a field whose name
matches the original type, the generated marshaling methods will use the overridden type
and convert to and from the original field type.

In this example, the specialString type implements json.Unmarshaler to enforce additional
parsing rules. When json.Unmarshal is used with type foo, the specialString unmarshaler
will be used to parse the value of SpecialField.

	//go:generate gencodec -dir . -type foo -field-override fooMarshaling -out foo_json.go

	type foo struct {
		Field        string
		SpecialField string
	}

	type fooMarshaling struct {
		SpecialField specialString // overrides type of SpecialField when marshaling/unmarshaling
	}

Relaxed Field Conversions

Field types in the override struct must be trivially convertible to the original field
type. gencodec's definition of 'convertible' is less restrictive than the usual rules
defined in the Go language specification.

The following conversions are supported:

If the fields are directly assignable, no conversion is generated.

If the fields are convertible according to Go language rules, a simple conversion is
emitted. Example input code:

	type specialString string

	func (s *specialString) UnmarshalJSON(input []byte) error { ... }

	type Foo struct{ S string }

	type fooMarshaling struct{ S specialString }

The generated code will contain:

	func (f *Foo) UnmarshalJSON(input []byte) error {
		var dec struct{ S *specialString }
		...
		f.S = string(dec.specialString)
		...
	}

If the fields are of map or slice type and the element (and key) types are convertible, a
simple loop is emitted. Example input code:

	type Foo2 struct{ M map[string]string }

	type foo2Marshaling struct{ S map[string]specialString }

The generated code is similar to this snippet:

	func (f *Foo2) UnmarshalJSON(input []byte) error {
		var dec struct{ M map[string]specialString }
		...
		for key, _ := range dec.M {
			f.M[key] = string(dec.M[key])
		}
		...
	}

*/
package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"io/ioutil"
	"os"
	"reflect"
	"strings"

	"golang.org/x/tools/imports"
)

func main() {
	var (
		pkgdir    = flag.String("dir", ".", "input package directory")
		output    = flag.String("out", "-", "output file")
		typename  = flag.String("type", "", "type to generate")
		overrides = flag.String("field-override", "", "type to take field type replacements from")
	)
	flag.Parse()

	cfg := Config{Dir: *pkgdir, Type: *typename, FieldOverride: *overrides}
	code, err := cfg.process()
	if err != nil {
		fatal(err)
	}
	if *output == "-" {
		os.Stdout.Write(code)
	} else if err := ioutil.WriteFile(*output, code, 0644); err != nil {
		fatal(err)
	}
}

func fatal(args ...interface{}) {
	fmt.Fprintln(os.Stderr, args...)
	os.Exit(1)
}

type Config struct {
	Dir           string // input package directory
	Type          string // type to generate methods for
	FieldOverride string // name of struct type for field overrides
	Importer      types.Importer
	FileSet       *token.FileSet
}

func (cfg *Config) process() (code []byte, err error) {
	if cfg.FileSet == nil {
		cfg.FileSet = token.NewFileSet()
	}
	if cfg.Importer == nil {
		cfg.Importer = importer.Default()
	}
	pkg, err := loadPackage(cfg)
	if err != nil {
		return nil, err
	}
	typ, err := lookupStructType(pkg.Scope(), cfg.Type)
	if err != nil {
		return nil, fmt.Errorf("can't find %s: %v", cfg.Type, err)
	}

	// Construct the marshaling type.
	mtyp := newMarshalerType(cfg.FileSet, cfg.Importer, typ)
	if cfg.FieldOverride != "" {
		otyp, err := lookupStructType(pkg.Scope(), cfg.FieldOverride)
		if err != nil {
			return nil, fmt.Errorf("can't find field replacement type %s: %v", cfg.FieldOverride, err)
		}
		err = mtyp.loadOverrides(cfg.FieldOverride, otyp.Underlying().(*types.Struct))
		if err != nil {
			return nil, err
		}
	}

	// Generate and format the output.
	// Use goimports to format the source because it separates imports.
	code = genPackage(mtyp)
	opt := &imports.Options{Comments: true, TabIndent: true, TabWidth: 8}
	code, err = imports.Process("", code, opt)
	if err != nil {
		return code, fmt.Errorf("can't gofmt generated code:", err, "\n"+string(code))
	}
	return code, nil
}

func loadPackage(cfg *Config) (*types.Package, error) {
	// Load the package.
	pkgs, err := parser.ParseDir(cfg.FileSet, cfg.Dir, nil, parser.AllErrors)
	if err != nil {
		return nil, err
	}
	if len(pkgs) == 0 || len(pkgs) > 1 {
		return nil, fmt.Errorf("input directory must contain exactly one package")
	}
	var files []*ast.File
	var name string
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			files = append(files, file)
		}
		name = pkg.Name
		break
	}
	// Type-check the package.
	tcfg := types.Config{IgnoreFuncBodies: true, FakeImportC: true, Importer: cfg.Importer}
	tpkg, err := tcfg.Check(name, cfg.FileSet, files, nil)
	if err != nil {
		return nil, err
	}
	return tpkg, nil
}

func genPackage(mtyp *marshalerType) []byte {
	w := new(bytes.Buffer)
	fmt.Fprintln(w, "// generated by gencodec, do not edit.\n")
	fmt.Fprintln(w, "package", mtyp.orig.Obj().Pkg().Name())
	fmt.Fprintln(w)
	mtyp.scope.writeImportDecl(w)
	fmt.Fprintln(w)
	writeFunction(w, mtyp.fs, genMarshalJSON(mtyp))
	fmt.Fprintln(w)
	writeFunction(w, mtyp.fs, genUnmarshalJSON(mtyp))
	fmt.Fprintln(w)
	writeFunction(w, mtyp.fs, genMarshalYAML(mtyp))
	fmt.Fprintln(w)
	writeFunction(w, mtyp.fs, genUnmarshalYAML(mtyp))
	return w.Bytes()
}

// marshalerType represents the intermediate struct type used during marshaling.
// This is the input data to all the Go code templates.
type marshalerType struct {
	name   string
	Fields []*marshalerField
	fs     *token.FileSet
	orig   *types.Named
	scope  *fileScope
}

// marshalerField represents a field of the intermediate marshaling type.
type marshalerField struct {
	name    string
	typ     types.Type
	origTyp types.Type
	tag     string
}

func newMarshalerType(fs *token.FileSet, imp types.Importer, typ *types.Named) *marshalerType {
	mtyp := &marshalerType{name: typ.Obj().Name(), fs: fs, orig: typ}
	styp := typ.Underlying().(*types.Struct)
	mtyp.scope = newFileScope(imp, typ.Obj().Pkg())
	mtyp.scope.addReferences(styp)

	// Add packages which are always needed.
	mtyp.scope.addImport("encoding/json")
	mtyp.scope.addImport("errors")

	for i := 0; i < styp.NumFields(); i++ {
		f := styp.Field(i)
		if !f.Exported() {
			continue
		}
		mf := &marshalerField{
			name:    f.Name(),
			typ:     ensurePointer(f.Type()),
			origTyp: f.Type(),
			tag:     styp.Tag(i),
		}
		if f.Anonymous() {
			fmt.Fprintf(os.Stderr, "Warning: ignoring embedded field %s\n", f.Name())
			continue
		}
		mtyp.Fields = append(mtyp.Fields, mf)
	}
	return mtyp
}

// loadOverrides sets field types of the intermediate marshaling type from
// matching fields of otyp.
func (mtyp *marshalerType) loadOverrides(otypename string, otyp *types.Struct) error {
	for i := 0; i < otyp.NumFields(); i++ {
		of := otyp.Field(i)
		if of.Anonymous() || !of.Exported() {
			return fmt.Errorf("%v: field override type cannot have embedded or unexported fields", mtyp.fs.Position(of.Pos()))
		}
		f := mtyp.fieldByName(of.Name())
		if f == nil {
			return fmt.Errorf("%v: no matching field for %s in original type %s", mtyp.fs.Position(of.Pos()), of.Name(), mtyp.name)
		}
		if err := checkConvertible(of.Type(), f.origTyp); err != nil {
			return fmt.Errorf("%v: invalid field override: %v", mtyp.fs.Position(of.Pos()), err)
		}
		f.typ = ensurePointer(of.Type())
	}
	mtyp.scope.addReferences(otyp)
	return nil
}

func (mtyp *marshalerType) fieldByName(name string) *marshalerField {
	for _, f := range mtyp.Fields {
		if f.name == name {
			return f
		}
	}
	return nil
}

// isOptional returns whether the field is optional when decoding the given format.
func (mf *marshalerField) isOptional(format string) bool {
	rtag := reflect.StructTag(mf.tag)
	if rtag.Get("optional") == "true" || rtag.Get("optional") == "yes" {
		return true
	}
	// Fields with json:"-" must be treated as optional.
	return strings.HasPrefix(rtag.Get(format), "-")
}

// encodedName returns the alternative field name assigned by the format's struct tag.
func (mf *marshalerField) encodedName(format string) string {
	val := reflect.StructTag(mf.tag).Get(format)
	if comma := strings.Index(val, ","); comma != -1 {
		val = val[:comma]
	}
	if val == "" || val == "-" {
		return uncapitalize(mf.name)
	}
	return val
}

func uncapitalize(s string) string {
	return strings.ToLower(s[:1]) + s[1:]
}
