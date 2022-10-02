// Copyright (c) 2022 DevDane <dane@danecwalker.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package lexer

import (
	"bufio"
	"io"
	"strings"
	"unicode"
)

type Position struct {
	Line   int
	Column int
}

type Lexer interface {
	Next() *Token
	Peek() rune
	Consume() rune

	IsEOF() bool
}

type lexer struct {
	input    *bufio.Reader
	position Position
}

func NewLexer(input *bufio.Reader) Lexer {
	return &lexer{
		input: input,
		position: Position{
			Line:   0,
			Column: 0,
		},
	}
}

func (l *lexer) IsEOF() bool {
	return l.Peek() == 0
}

func (l *lexer) Peek() rune {
	b, err := l.input.Peek(1)
	if err != nil {
		if err == io.EOF {
			return 0
		}
		panic(err)
	}
	return rune(b[0])
}

func (l *lexer) Consume() rune {
	r, _, err := l.input.ReadRune()
	if err != nil {
		if err == io.EOF {
			return 0
		}
		panic(err)
	}
	l.position.Column++
	return r
}

func (l *lexer) ConsumeWhile(predicate func(rune) bool) string {
	var sb strings.Builder
	for {
		ch := l.Peek()
		if !predicate(ch) {
			break
		}
		sb.WriteRune(l.Consume())
	}
	return sb.String()
}

func (l *lexer) ConsumeWhitespace() {
	l.ConsumeWhile(func(ch rune) bool {
		if unicode.IsSpace(ch) {
			if ch == '\n' {
				l.position.Line++
				l.position.Column = 1
			}
			return true
		} else {
			return false
		}
	})
}

func (l *lexer) Next() *Token {
	l.ConsumeWhitespace()
	switch l.Peek() {
	case 0:
		return l.newToken(EOF, string(l.Consume()))
	case '=':
		lit := string(l.Consume())
		if l.Peek() == '=' {
			lit += string(l.Consume())
			return l.newToken(EQ, lit)
		}
		return l.newToken(ASSIGN, string(lit))
	case '+':
		lit := string(l.Consume())
		if l.Peek() == '=' {
			lit += string(l.Consume())
			return l.newToken(PLUS_ASSIGN, lit)
		}
		return l.newToken(PLUS, lit)
	case '-':
		lit := string(l.Consume())
		if l.Peek() == '=' {
			lit += string(l.Consume())
			return l.newToken(MINUS_ASSIGN, lit)
		}
		return l.newToken(MINUS, lit)
	case '!':
		lit := string(l.Consume())
		if l.Peek() == '=' {
			lit += string(l.Consume())
			return l.newToken(NOT_EQ, lit)
		}
		return l.newToken(BANG, string(lit))
	case '*':
		lit := string(l.Consume())
		if l.Peek() == '=' {
			lit += string(l.Consume())
			return l.newToken(ASTER_ASSIGN, lit)
		}
		return l.newToken(ASTER, string(l.Consume()))
	case '%':
		lit := string(l.Consume())
		if l.Peek() == '=' {
			lit += string(l.Consume())
			return l.newToken(MOD_ASSIGN, lit)
		}
		return l.newToken(MOD, string(l.Consume()))
	case '/':
		lit := string(l.Consume())
		if l.Peek() == '=' {
			lit += string(l.Consume())
			return l.newToken(SLASH_ASSIGN, lit)
		}
		return l.newToken(SLASH, string(l.Consume()))
	case '<':
		lit := string(l.Consume())
		if l.Peek() == '=' {
			lit += string(l.Consume())
			return l.newToken(LT_EQ, lit)
		}
		return l.newToken(LT, lit)
	case '>':
		lit := string(l.Consume())
		if l.Peek() == '=' {
			lit += string(l.Consume())
			return l.newToken(GT_EQ, lit)
		}
		return l.newToken(GT, lit)
	case ',':
		return l.newToken(COMMA, string(l.Consume()))
	case ';':
		return l.newToken(SEMICOLON, string(l.Consume()))
	case '(':
		return l.newToken(LPAREN, string(l.Consume()))
	case ')':
		return l.newToken(RPAREN, string(l.Consume()))
	case '{':
		return l.newToken(LBRACE, string(l.Consume()))
	case '}':
		return l.newToken(RBRACE, string(l.Consume()))
	case '"':
		l.Consume()
		lit := l.ConsumeWhile(func(r rune) bool { return r != '"' })
		l.Consume()
		return l.newToken(STRING, lit)
	default:
		if unicode.IsLetter(l.Peek()) {
			lit := l.ConsumeWhile(unicode.IsLetter)
			return l.newToken(lookupIdent(lit), lit)
		} else if unicode.IsDigit(l.Peek()) {
			lit := l.ConsumeWhile(unicode.IsDigit)
			return l.newToken(INT, lit)
		} else {
			return l.newToken(ILLEGAL, string(l.Consume()))
		}
	}
}

func lookupIdent(ident string) TokenType {
	keyword, ok := keywords[ident]
	if ok {
		return keyword
	} else {
		return IDENT
	}
}
