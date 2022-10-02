// Copyright (c) 2022 DevDane <dane@danecwalker.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package object

import (
	"fmt"

	"github.com/danecwalker/ponic/engine/ast"
)

type Object interface {
	Type() Type
	Inspect() string
	String() string
}

type Type int

const (
	INTEGER Type = iota
	BOOLEAN
	NULL
	STRING
	FUNCTION
)

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() Type {
	return rv.Value.Type()
}

func (rv *ReturnValue) Inspect() string {
	return rv.Value.Inspect()
}

func (rv *ReturnValue) String() string {
	return rv.Value.String()
}

type Builtin struct {
	Func func(args ...Object) Object
}

func (b *Builtin) Type() Type {
	return FUNCTION
}
func (b *Builtin) Inspect() string {
	return "builtin function"
}
func (b *Builtin) String() string {
	return b.Inspect()
}

type Integer struct {
	Value int64
}

func (i *Integer) Type() Type {
	return INTEGER
}
func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}
func (i *Integer) String() string {
	return fmt.Sprintf("Int(%d)", i.Value)
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() Type {
	return BOOLEAN
}
func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}
func (b *Boolean) String() string {
	return fmt.Sprintf("Boolean(%t)", b.Value)
}

type Null struct{}

func (n *Null) Type() Type {
	return NULL
}
func (n *Null) Inspect() string {
	return "null"
}
func (n *Null) String() string {
	return "Null()"
}

type String struct {
	Value string
}

func (s *String) Type() Type {
	return STRING
}
func (s *String) Inspect() string {
	return s.Value
}
func (s *String) String() string {
	return fmt.Sprintf("String(%s)", s.Value)
}

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Scope      *Scope
}

func (f *Function) Type() Type {
	return FUNCTION
}
func (f *Function) Inspect() string {
	return "function"
}
func (f *Function) String() string {
	return "Function()"
}
