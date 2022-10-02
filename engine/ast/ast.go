// Copyright (c) 2022 DevDane <dane@danecwalker.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package ast

import (
	"fmt"

	"github.com/danecwalker/ponic/engine/lexer"
)

type Node interface {
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type AST struct {
	Statements []Statement
}

func (a *AST) String() string {
	return fmt.Sprintf("AST(%s)", a.Statements)
}

type Identifier struct {
	Token *lexer.Token
	Value string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) String() string {
	return fmt.Sprintf("Identifier(%s)", i.Value)
}

type ExpressionStatement struct {
	Token      *lexer.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) String() string {
	return es.Expression.String()
}

type LetStatement struct {
	Token *lexer.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) String() string {
	return fmt.Sprintf("LetStatement(%s, %s)", ls.Name, ls.Value)
}

type ConstStatement struct {
	Token *lexer.Token
	Name  *Identifier
	Value Expression
}

func (cs *ConstStatement) statementNode() {}
func (cs *ConstStatement) String() string {
	return fmt.Sprintf("ConstStatement(%s, %s)", cs.Name, cs.Value)
}

type ReturnStatement struct {
	Token       *lexer.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) String() string {
	return fmt.Sprintf("ReturnStatement(%s)", rs.ReturnValue)
}

type UnOp struct {
	Token    *lexer.Token
	Operator string
	Right    Expression
}

func (uo *UnOp) expressionNode() {}
func (uo *UnOp) String() string {
	return fmt.Sprintf("UnOp(%s, %s)", uo.Operator, uo.Right)
}

type BinOp struct {
	Token    *lexer.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (bo *BinOp) expressionNode() {}
func (bo *BinOp) String() string {
	return fmt.Sprintf("BinOp(%s, %s, %s)", bo.Left, bo.Operator, bo.Right)
}

type StringLiteral struct {
	Token *lexer.Token
	Value string
}

func (sl *StringLiteral) expressionNode() {}
func (sl *StringLiteral) String() string {
	return fmt.Sprintf("String(%s)", sl.Value)
}

type IntegerLiteral struct {
	Token *lexer.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode() {}
func (il *IntegerLiteral) String() string {
	return fmt.Sprintf("Int(%d)", il.Value)
}

type BooleanLiteral struct {
	Token *lexer.Token
	Value bool
}

func (b *BooleanLiteral) expressionNode() {}
func (b *BooleanLiteral) String() string {
	return fmt.Sprintf("Bool(%t)", b.Value)
}

type IfExpression struct {
	Token       *lexer.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode() {}
func (ie *IfExpression) String() string {
	return fmt.Sprintf("IfExpression(%s, %s, %s)", ie.Condition, ie.Consequence, ie.Alternative)
}

type BlockStatement struct {
	Token      *lexer.Token
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}
func (bs *BlockStatement) String() string {
	return fmt.Sprintf("BlockStatement(%s)", bs.Statements)
}

type FunctionLiteral struct {
	Token      *lexer.Token
	Parameters []*Identifier
	Body       *BlockStatement
	Name       *Identifier
	Named      bool
}

func (fl *FunctionLiteral) expressionNode() {}
func (fl *FunctionLiteral) String() string {
	if fl.Named {
		return fmt.Sprintf("NamedFunctionLiteral(%s, %s, %s)", fl.Name, fl.Parameters, fl.Body)
	}
	return fmt.Sprintf("FunctionLiteral(%s, %s)", fl.Parameters, fl.Body)
}

type CallExpression struct {
	Token     *lexer.Token
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode() {}
func (ce *CallExpression) String() string {
	return fmt.Sprintf("CallExpression(%s, %s)", ce.Function, ce.Arguments)
}
