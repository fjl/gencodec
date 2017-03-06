// Copyright 2017 Felix Lange <fjl@twurst.com>.
// Use of this source code is governed by the MIT license,
// which can be found in the LICENSE file.

package main

import (
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"go/types"
	"io"
	"strconv"
	"strings"

	. "github.com/garslo/gogen"
)

var NIL = Name("nil")

type marshalMethod struct {
	mtyp  *marshalerType
	scope *funcScope
}

func newMarshalMethod(mtyp *marshalerType) *marshalMethod {
	return &marshalMethod{mtyp, newFuncScope(mtyp.scope)}
}

func writeFunction(w io.Writer, fs *token.FileSet, fn Function) {
	printer.Fprint(w, fs, fn.Declaration())
	fmt.Fprintln(w)
}

// genUnmarshalJSON generates the UnmarshalJSON method.
func genUnmarshalJSON(mtyp *marshalerType) Function {
	var (
		m        = newMarshalMethod(mtyp)
		recv     = m.receiver()
		input    = Name(m.scope.newIdent("input"))
		intertyp = m.intermediateType(m.scope.newIdent(m.mtyp.orig.Obj().Name() + "JSON"))
		dec      = Name(m.scope.newIdent("dec"))
		conv     = Name(m.scope.newIdent("x"))
		json     = Name(m.scope.parent.packageName("encoding/json"))
	)
	fn := Function{
		Receiver:    recv,
		Name:        "UnmarshalJSON",
		ReturnTypes: Types{{TypeName: "error"}},
		Parameters:  Types{{Name: input.Name, TypeName: "[]byte"}},
		Body: []Statement{
			declStmt{intertyp},
			Declare{Name: dec.Name, TypeName: intertyp.Name},
			errCheck(CallFunction{
				Func:   Dotted{Receiver: json, Name: "Unmarshal"},
				Params: []Expression{input, AddressOf{Value: dec}},
			}),
			Declare{Name: conv.Name, TypeName: m.mtyp.name},
		},
	}
	fn.Body = append(fn.Body, m.unmarshalConversions(dec, conv, "json")...)
	fn.Body = append(fn.Body, Assign{Lhs: Star{Value: Name(recv.Name)}, Rhs: conv})
	fn.Body = append(fn.Body, Return{Values: []Expression{Name("nil")}})
	return fn
}

// genMarshalJSON generates the MarshalJSON method.
func genMarshalJSON(mtyp *marshalerType) Function {
	var (
		m        = newMarshalMethod(mtyp)
		recv     = m.receiver()
		intertyp = m.intermediateType(m.scope.newIdent(m.mtyp.orig.Obj().Name() + "JSON"))
		enc      = Name(m.scope.newIdent("enc"))
		json     = Name(m.scope.parent.packageName("encoding/json"))
	)
	fn := Function{
		Receiver:    recv,
		Name:        "MarshalJSON",
		ReturnTypes: Types{{TypeName: "[]byte"}, {TypeName: "error"}},
		Body: []Statement{
			declStmt{intertyp},
			Declare{Name: enc.Name, TypeName: intertyp.Name},
		},
	}
	fn.Body = append(fn.Body, m.marshalConversions(Name(recv.Name), enc, "json")...)
	fn.Body = append(fn.Body, Return{Values: []Expression{
		CallFunction{
			Func:   Dotted{Receiver: json, Name: "Marshal"},
			Params: []Expression{AddressOf{Value: enc}},
		},
	}})
	return fn
}

// genUnmarshalYAML generates the UnmarshalYAML method.
func genUnmarshalYAML(mtyp *marshalerType) Function {
	var (
		m         = newMarshalMethod(mtyp)
		recv      = m.receiver()
		unmarshal = Name(m.scope.newIdent("unmarshal"))
		intertyp  = m.intermediateType(m.scope.newIdent(m.mtyp.orig.Obj().Name() + "YAML"))
		dec       = Name(m.scope.newIdent("dec"))
		conv      = Name(m.scope.newIdent("x"))
	)
	fn := Function{
		Receiver:    recv,
		Name:        "UnmarshalYAML",
		ReturnTypes: Types{{TypeName: "error"}},
		Parameters:  Types{{Name: unmarshal.Name, TypeName: "func (interface{}) error"}},
		Body: []Statement{
			declStmt{intertyp},
			Declare{Name: dec.Name, TypeName: intertyp.Name},
			errCheck(CallFunction{Func: unmarshal, Params: []Expression{AddressOf{Value: dec}}}),
			Declare{Name: conv.Name, TypeName: m.mtyp.name},
		},
	}
	fn.Body = append(fn.Body, m.unmarshalConversions(dec, conv, "yaml")...)
	fn.Body = append(fn.Body, Assign{Lhs: Star{Value: Name(recv.Name)}, Rhs: conv})
	fn.Body = append(fn.Body, Return{Values: []Expression{Name("nil")}})
	return fn
}

// genMarshalYAML generates the MarshalYAML method.
func genMarshalYAML(mtyp *marshalerType) Function {
	var (
		m        = newMarshalMethod(mtyp)
		recv     = m.receiver()
		intertyp = m.intermediateType(m.scope.newIdent(m.mtyp.orig.Obj().Name() + "YAML"))
		enc      = Name(m.scope.newIdent("enc"))
	)
	fn := Function{
		Receiver:    recv,
		Name:        "MarshalYAML",
		ReturnTypes: Types{{TypeName: "interface{}"}, {TypeName: "error"}},
		Body: []Statement{
			declStmt{intertyp},
			Declare{Name: enc.Name, TypeName: intertyp.Name},
		},
	}
	fn.Body = append(fn.Body, m.marshalConversions(Name(recv.Name), enc, "yaml")...)
	fn.Body = append(fn.Body, Return{Values: []Expression{AddressOf{Value: enc}}})
	return fn
}

func (m *marshalMethod) receiver() Receiver {
	letter := strings.ToLower(m.mtyp.name[:1])
	return Receiver{Name: m.scope.newIdent(letter), Type: Star{Value: Name(m.mtyp.name)}}
}

func (m *marshalMethod) intermediateType(name string) Struct {
	s := Struct{Name: name}
	for _, f := range m.mtyp.Fields {
		s.Fields = append(s.Fields, Field{
			Name:     f.name,
			TypeName: types.TypeString(f.typ, m.mtyp.scope.qualify),
			Tag:      f.tag,
		})
	}
	return s
}

func (m *marshalMethod) unmarshalConversions(from, to Var, format string) (s []Statement) {
	for _, f := range m.mtyp.Fields {
		accessFrom := Dotted{Receiver: from, Name: f.name}
		accessTo := Dotted{Receiver: to, Name: f.name}
		if f.isOptional(format) {
			s = append(s, If{
				Condition: NotEqual{Lhs: accessFrom, Rhs: NIL},
				Body:      m.convert(accessFrom, accessTo, f.typ, f.origTyp),
			})
		} else {
			err := fmt.Sprintf("missing required field %s for %s", f.encodedName(format), m.mtyp.name)
			errors := m.scope.parent.packageName("errors")
			s = append(s, If{
				Condition: Equals{Lhs: accessFrom, Rhs: NIL},
				Body: []Statement{
					Return{
						Values: []Expression{
							CallFunction{
								Func:   Dotted{Receiver: Name(errors), Name: "New"},
								Params: []Expression{stringLit{err}},
							},
						},
					},
				},
			})
			s = append(s, m.convert(accessFrom, accessTo, f.typ, f.origTyp)...)
		}
	}
	return s
}

func (m *marshalMethod) marshalConversions(from, to Var, format string) (s []Statement) {
	for _, f := range m.mtyp.Fields {
		accessFrom := Dotted{Receiver: from, Name: f.name}
		accessTo := Dotted{Receiver: to, Name: f.name}
		s = append(s, m.convert(accessFrom, accessTo, f.origTyp, f.typ)...)
	}
	return s
}

func (m *marshalMethod) convert(from, to Expression, fromtyp, totyp types.Type) (s []Statement) {
	// Remove pointer introduced by ensurePointer during field building.
	if isPointer(fromtyp) && !isPointer(totyp) {
		tmp := Name(m.scope.newIdent("v"))
		s = append(s, DeclareAndAssign{Lhs: tmp, Rhs: Star{Value: from}})
		from = tmp
		fromtyp = fromtyp.(*types.Pointer).Elem()
	} else if !isPointer(fromtyp) && isPointer(totyp) {
		tmp := Name(m.scope.newIdent("v"))
		s = append(s, DeclareAndAssign{Lhs: tmp, Rhs: AddressOf{Value: from}})
		from = tmp
		fromtyp = types.NewPointer(fromtyp)
	}

	// Generate the conversion.
	qf := m.mtyp.scope.qualify
	switch {
	case types.ConvertibleTo(fromtyp, totyp):
		s = append(s, Assign{Lhs: to, Rhs: simpleConv(from, fromtyp, totyp, qf)})
	case isSlice(fromtyp):
		fromElem := fromtyp.(*types.Slice).Elem()
		toElem := totyp.(*types.Slice).Elem()
		key := Name(m.scope.newIdent("i"))
		s = append(s, Assign{Lhs: to, Rhs: makeExpr(totyp, from, qf)})
		s = append(s, Range{
			Key: key, RangeValue: from,
			Body: []Statement{Assign{
				Lhs: Index{Value: to, Index: key},
				Rhs: simpleConv(Index{Value: from, Index: key}, fromElem, toElem, qf),
			}},
		})
	case isMap(fromtyp):
		fromKey, fromElem := fromtyp.(*types.Map).Key(), fromtyp.(*types.Map).Elem()
		toKey, toElem := totyp.(*types.Map).Key(), totyp.(*types.Map).Elem()
		key := Name(m.scope.newIdent("key"))
		s = append(s, Assign{Lhs: to, Rhs: makeExpr(totyp, from, qf)})
		s = append(s, Range{
			Key: key, RangeValue: from,
			Body: []Statement{Assign{
				Lhs: Index{Value: to, Index: simpleConv(key, fromKey, toKey, qf)},
				Rhs: simpleConv(Index{Value: from, Index: key}, fromElem, toElem, qf),
			}},
		})
	default:
		invalidConv(fromtyp, totyp, qf)
	}
	return s
}

func simpleConv(from Expression, fromtyp, totyp types.Type, qf types.Qualifier) Expression {
	if types.AssignableTo(fromtyp, totyp) {
		return from
	}
	if !types.ConvertibleTo(fromtyp, totyp) {
		invalidConv(fromtyp, totyp, qf)
	}
	toname := types.TypeString(totyp, qf)
	if isPointer(totyp) {
		toname = "(" + toname + ")" // hack alert!
	}
	return CallFunction{Func: Name(toname), Params: []Expression{from}}
}

func invalidConv(from, to types.Type, qf types.Qualifier) {
	panic(fmt.Errorf("BUG: invalid conversion %s -> %s", types.TypeString(from, qf), types.TypeString(to, qf)))
}

func makeExpr(typ types.Type, lenfrom Expression, qf types.Qualifier) Expression {
	return CallFunction{Func: Name("make"), Params: []Expression{
		Name(types.TypeString(typ, qf)),
		CallFunction{Func: Name("len"), Params: []Expression{lenfrom}},
	}}
}

func errCheck(expr Expression) If {
	err := Name("err")
	return If{
		Init:      DeclareAndAssign{Lhs: err, Rhs: expr},
		Condition: Equals{Lhs: err, Rhs: NIL},
		Body:      []Statement{Return{Values: []Expression{err}}},
	}
}

type stringLit struct {
	V string
}

func (l stringLit) Expression() ast.Expr {
	return &ast.BasicLit{Kind: token.STRING, Value: strconv.Quote(l.V)}
}

type declStmt struct {
	d Declaration
}

func (ds declStmt) Statement() ast.Stmt {
	return &ast.DeclStmt{Decl: ds.d.Declaration()}
}
