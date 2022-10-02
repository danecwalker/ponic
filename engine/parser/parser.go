// Copyright (c) 2022 DevDane <dane@danecwalker.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package parser

import (
	"strconv"

	"github.com/danecwalker/ponic/engine/ast"
	"github.com/danecwalker/ponic/engine/lexer"
)

type Parser interface {
	Parse() *ast.AST
}

type (
	nud func() ast.Expression
	led func(ast.Expression) ast.Expression
)

type parser struct {
	lexer     lexer.Lexer
	curToken  *lexer.Token
	peekToken *lexer.Token
	nuds      map[lexer.TokenType]nud
	leds      map[lexer.TokenType]led
}

func NewParser(l lexer.Lexer) Parser {
	p := &parser{
		lexer:     l,
		curToken:  nil,
		peekToken: l.Next(),
		nuds:      make(map[lexer.TokenType]nud),
		leds:      make(map[lexer.TokenType]led),
	}

	p.registerNud(lexer.IDENT, p.parseIdentifier)
	p.registerNud(lexer.INT, p.parseIntegerLiteral)
	p.registerNud(lexer.TRUE, p.parseBooleanLiteral)
	p.registerNud(lexer.FALSE, p.parseBooleanLiteral)
	p.registerNud(lexer.STRING, p.parseStringLiteral)

	p.registerNud(lexer.FUNCTION, p.parseFunctionLiteral)
	p.registerNud(lexer.IF, p.parseIfExpression)
	p.registerNud(lexer.FOR, p.parseForExpression)

	p.registerNud(lexer.MINUS, p.parsePrefixExpression)
	p.registerNud(lexer.BANG, p.parsePrefixExpression)
	p.registerNud(lexer.LPAREN, p.parseGroupedExpression)

	p.registerLed(lexer.LPAREN, p.parseCallExpression)

	p.registerLed(lexer.ASSIGN, p.parseInfixExpression)
	p.registerLed(lexer.PLUS_ASSIGN, p.parseInfixExpression)
	p.registerLed(lexer.MINUS_ASSIGN, p.parseInfixExpression)
	p.registerLed(lexer.ASTER_ASSIGN, p.parseInfixExpression)
	p.registerLed(lexer.SLASH_ASSIGN, p.parseInfixExpression)
	p.registerLed(lexer.MOD_ASSIGN, p.parseInfixExpression)

	p.registerLed(lexer.PLUS, p.parseInfixExpression)
	p.registerLed(lexer.MINUS, p.parseInfixExpression)
	p.registerLed(lexer.ASTER, p.parseInfixExpression)
	p.registerLed(lexer.SLASH, p.parseInfixExpression)
	p.registerLed(lexer.MOD, p.parseInfixExpression)
	p.registerLed(lexer.LT, p.parseInfixExpression)
	p.registerLed(lexer.GT, p.parseInfixExpression)
	p.registerLed(lexer.EQ, p.parseInfixExpression)
	p.registerLed(lexer.NOT_EQ, p.parseInfixExpression)
	p.registerLed(lexer.LT_EQ, p.parseInfixExpression)
	p.registerLed(lexer.GT_EQ, p.parseInfixExpression)

	return p
}

func (p *parser) registerNud(t lexer.TokenType, f nud) {
	p.nuds[t] = f
}

func (p *parser) registerLed(t lexer.TokenType, f led) {
	p.leds[t] = f
}

func (p *parser) next() *lexer.Token {
	return p.peekToken
}

func (p *parser) eat() *lexer.Token {
	p.curToken = p.peekToken
	p.peekToken = p.lexer.Next()
	return p.curToken
}

func (p *parser) isNext(t lexer.TokenType) bool {
	return p.next().Type == t
}

func (p *parser) Parse() *ast.AST {
	program := &ast.AST{}
	program.Statements = []ast.Statement{}

	for !p.isNext(lexer.EOF) {
		stmt := p.parseStatement()
		program.Statements = append(program.Statements, stmt)
	}

	return program
}

func (p *parser) parseStatement() ast.Statement {
	switch p.next().Type {
	case lexer.LET:
		return p.parseLetStatement()
	case lexer.CONST:
		return p.parseConstStatement()
	case lexer.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.eat()}

	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.isNext(lexer.SEMICOLON) {
		p.eat()
	}

	return stmt
}

func (p *parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.next()}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.isNext(lexer.SEMICOLON) {
		p.eat()
	}

	return stmt
}

func (p *parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.eat()}

	if !p.isNext(lexer.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.eat(), Value: p.curToken.Literal}

	if !p.isNext(lexer.ASSIGN) {
		return nil
	}

	p.eat()

	stmt.Value = p.parseExpression(LOWEST)

	if p.isNext(lexer.SEMICOLON) {
		p.eat()
	}

	return stmt
}

func (p *parser) parseConstStatement() *ast.ConstStatement {
	stmt := &ast.ConstStatement{Token: p.eat()}

	if !p.isNext(lexer.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.eat(), Value: p.curToken.Literal}

	if !p.isNext(lexer.ASSIGN) {
		return nil
	}

	p.eat()

	stmt.Value = p.parseExpression(LOWEST)

	if p.isNext(lexer.SEMICOLON) {
		p.eat()
	}

	return stmt
}

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	ASSIGN      // =
	CALL        // myFunction(X)
)

var precedences = map[lexer.TokenType]int{
	lexer.EQ:           EQUALS,
	lexer.NOT_EQ:       EQUALS,
	lexer.LT_EQ:        EQUALS,
	lexer.GT_EQ:        EQUALS,
	lexer.LT:           LESSGREATER,
	lexer.GT:           LESSGREATER,
	lexer.PLUS:         SUM,
	lexer.MINUS:        SUM,
	lexer.SLASH:        PRODUCT,
	lexer.ASTER:        PRODUCT,
	lexer.MOD:          PRODUCT,
	lexer.ASSIGN:       ASSIGN,
	lexer.PLUS_ASSIGN:  ASSIGN,
	lexer.MINUS_ASSIGN: ASSIGN,
	lexer.ASTER_ASSIGN: ASSIGN,
	lexer.SLASH_ASSIGN: ASSIGN,
	lexer.MOD_ASSIGN:   ASSIGN,
	lexer.LPAREN:       CALL,
}

func (p *parser) nextPrecedence() int {
	if p, ok := precedences[p.next().Type]; ok {
		return p
	}

	return LOWEST
}

func (p *parser) parseExpression(precedence int) ast.Expression {
	p.eat()
	_nud := p.nuds[p.curToken.Type]
	if _nud == nil {
		return nil
	}
	left := _nud()

	for !p.isNext(lexer.SEMICOLON) && precedence < p.nextPrecedence() {
		_led := p.leds[p.next().Type]
		if _led == nil {
			return left
		}

		p.eat()
		left = _led(left)
	}

	return left
}

func (p *parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *parser) parseIntegerLiteral() ast.Expression {
	il, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		panic(err)
	}

	return &ast.IntegerLiteral{Token: p.curToken, Value: il}
}

func (p *parser) parseBooleanLiteral() ast.Expression {
	return &ast.BooleanLiteral{Token: p.curToken, Value: p.curToken.Type == lexer.TRUE}
}

func (p *parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *parser) parseGroupedExpression() ast.Expression {
	exp := p.parseExpression(LOWEST)

	if !p.isNext(lexer.RPAREN) {
		return nil
	}

	p.eat()

	return exp
}

func (p *parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expr := &ast.BinOp{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := precedences[p.curToken.Type]
	expr.Right = p.parseExpression(precedence)

	return expr
}

func (p *parser) parsePrefixExpression() ast.Expression {
	expr := &ast.UnOp{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	precedence := precedences[p.curToken.Type]
	expr.Right = p.parseExpression(precedence)

	return expr
}

func (p *parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken, Named: false}

	if p.isNext(lexer.IDENT) {
		lit.Name = &ast.Identifier{Token: p.eat(), Value: p.curToken.Literal}
		lit.Named = true
	}

	if !p.isNext(lexer.LPAREN) {
		return nil
	}

	p.eat()

	lit.Parameters = p.parseFunctionParams()

	if !p.isNext(lexer.LBRACE) {
		return nil
	}
	p.eat()

	lit.Body = p.parseBlockStatement()

	return lit
}

func (p *parser) parseFunctionParams() []*ast.Identifier {
	idents := []*ast.Identifier{}

	if p.isNext(lexer.RPAREN) {
		p.eat()
		return idents
	}

	ident := &ast.Identifier{Token: p.eat(), Value: p.curToken.Literal}
	idents = append(idents, ident)

	for p.isNext(lexer.COMMA) {
		p.eat()
		ident := &ast.Identifier{Token: p.eat(), Value: p.curToken.Literal}
		idents = append(idents, ident)
	}

	if !p.isNext(lexer.RPAREN) {
		return nil
	}
	p.eat()

	return idents
}

func (p *parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}

	for !p.isNext(lexer.RBRACE) && !p.isNext(lexer.EOF) {
		stmt := p.parseStatement()
		block.Statements = append(block.Statements, stmt)
	}

	return block
}

func (p *parser) parseIfExpression() ast.Expression {
	exp := &ast.IfExpression{Token: p.curToken}

	if !p.isNext(lexer.LPAREN) {
		return nil
	}
	p.eat()

	exp.Condition = p.parseExpression(LOWEST)

	if !p.isNext(lexer.RPAREN) {
		return nil
	}
	p.eat()

	if !p.isNext(lexer.LBRACE) {
		return nil
	}
	p.eat()

	exp.Consequence = p.parseBlockStatement()

	if !p.isNext(lexer.RBRACE) {
		return nil
	}
	p.eat()

	if p.isNext(lexer.ELSE) {
		p.eat()

		if !p.isNext(lexer.LBRACE) {
			return nil
		}
		p.eat()

		exp.Alternative = p.parseBlockStatement()

		if !p.isNext(lexer.RBRACE) {
			return nil
		}
		p.eat()
	}

	return exp
}

func (p *parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseCallArguments()
	return exp
}

func (p *parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	if p.isNext(lexer.RPAREN) {
		p.eat()
		return args
	}

	args = append(args, p.parseExpression(LOWEST))

	for p.isNext(lexer.COMMA) {
		p.eat()
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.isNext(lexer.RPAREN) {
		return nil
	}
	p.eat()

	return args
}

func (p *parser) parseForExpression() ast.Expression {
	exp := &ast.ForExpression{Token: p.curToken}

	if !p.isNext(lexer.LPAREN) {
		return nil
	}
	p.eat()

	if !p.isNext(lexer.LET) {
		return nil
	}

	exp.Initializer = p.parseLetStatement()

	exp.Condition = p.parseExpression(LOWEST)

	if !p.isNext(lexer.SEMICOLON) {
		return nil
	}
	p.eat()

	exp.Post = p.parseExpressionStatement()

	if !p.isNext(lexer.RPAREN) {
		return nil
	}
	p.eat()

	if !p.isNext(lexer.LBRACE) {
		return nil
	}
	p.eat()

	exp.Body = p.parseBlockStatement()

	if !p.isNext(lexer.RBRACE) {
		return nil
	}
	p.eat()

	return exp
}
