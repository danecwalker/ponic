// Copyright (c) 2022 DevDane <dane@danecwalker.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package lexer

import "fmt"

type Token struct {
	Type    TokenType
	Literal string
	Pos     Position
}

func (t Token) String() string {
	return fmt.Sprintf("Token(%s, %s, %d:%d)", t.Type, t.Literal, t.Pos.Line, t.Pos.Column)
}

func (l *lexer) newToken(tokenType TokenType, lit string) *Token {
	return &Token{
		Type:    tokenType,
		Literal: lit,
		Pos:     l.position,
	}
}

type TokenType int

const (
	// Special tokens
	ILLEGAL TokenType = iota
	EOF

	// Identifiers + literals
	IDENT  // main
	INT    // 12345
	STRING // "hello world"

	// Operators
	ASSIGN       // =
	PLUS         // +
	MINUS        // -
	BANG         // !
	ASTER        // *
	SLASH        // /
	MOD          // %
	PLUS_ASSIGN  // +=
	MINUS_ASSIGN // -=
	ASTER_ASSIGN // *=
	SLASH_ASSIGN // /=
	MOD_ASSIGN   // %=

	LT // <
	GT // >

	EQ     // ==
	NOT_EQ // !=
	LT_EQ  // <=
	GT_EQ  // >=

	// Delimiters
	COMMA     // ,
	SEMICOLON // ;

	LPAREN // (
	RPAREN // )
	LBRACE // {
	RBRACE // }

	// Keywords
	FUNCTION // fn
	LET      // let
	CONST    // const
	TRUE     // true
	FALSE    // false
	IF       // if
	ELSE     // else
	FOR      // for
	RETURN   // return
)

var TokenMap = [...]string{
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",

	IDENT:  "IDENT",
	INT:    "INT",
	STRING: "STRING",

	ASSIGN:       "=",
	PLUS:         "+",
	MINUS:        "-",
	BANG:         "!",
	ASTER:        "*",
	SLASH:        "/",
	MOD:          "%",
	PLUS_ASSIGN:  "+=",
	MINUS_ASSIGN: "-=",
	ASTER_ASSIGN: "*=",
	SLASH_ASSIGN: "/=",
	MOD_ASSIGN:   "%=",

	LT: "<",
	GT: ">",

	EQ:     "==",
	NOT_EQ: "!=",
	LT_EQ:  "<=",
	GT_EQ:  ">=",

	COMMA:     ",",
	SEMICOLON: ";",

	LPAREN: "(",
	RPAREN: ")",
	LBRACE: "{",
	RBRACE: "}",

	FUNCTION: "FUNCTION",
	LET:      "LET",
	CONST:    "CONST",
	TRUE:     "TRUE",
	FALSE:    "FALSE",
	IF:       "IF",
	ELSE:     "ELSE",
	FOR:      "FOR",
	RETURN:   "RETURN",
}

func (t TokenType) String() string {
	return TokenMap[t]
}

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"const":  CONST,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"for":    FOR,
	"return": RETURN,
}
