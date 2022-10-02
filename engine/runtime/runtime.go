// Copyright (c) 2022 DevDane <dane@danecwalker.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT
package runtime

import (
	"github.com/danecwalker/ponic/engine/ast"
	"github.com/danecwalker/ponic/engine/object"
)

func Run(node ast.Node, scope *object.Scope) object.Object {
	switch node := (node).(type) {
	case *ast.AST:
		return runAST(node.Statements, scope)
	case *ast.ExpressionStatement:
		return Run(node.Expression, scope)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.BooleanLiteral:
		return &object.Boolean{Value: node.Value}
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.LetStatement:
		val := Run(node.Value, scope)
		scope.Set(node.Name.Value, val, object.LET)
	case *ast.ConstStatement:
		val := Run(node.Value, scope)
		scope.Set(node.Name.Value, val, object.CONST)
	case *ast.ForExpression:
		return runForExpression(node, scope)
	case *ast.ReturnStatement:
		val := Run(node.ReturnValue, scope)
		return &object.ReturnValue{Value: val}
	case *ast.Identifier:
		val, ok := scope.Get(node.Value)
		if !ok {
			if findBuiltin, ok := Builtins[node.Value]; ok {
				return &object.Builtin{Func: findBuiltin}
			}
			panic("Undefined variable " + node.Value)
		}
		return val
	case *ast.UnOp:
		right := Run(node.Right, scope)
		return runUnop(node.Operator, right)
	case *ast.BinOp:
		if node.Operator == "=" || node.Operator == "+=" || node.Operator == "-=" || node.Operator == "*=" || node.Operator == "/=" || node.Operator == "%=" {
			switch n := node.Left.(type) {
			case *ast.Identifier:
				return runRebind(n, node.Right, node.Operator, scope)
			}
		}
		left := Run(node.Left, scope)
		right := Run(node.Right, scope)
		return runBinop(node.Operator, left, right)
	case *ast.IfExpression:
		return runIfExpression(node, scope)
	case *ast.BlockStatement:
		return runBlockStatement(node, scope)
	case *ast.FunctionLiteral:
		s := object.NewScope()
		s.Parent = scope
		if node.Named {
			scope.Set(node.Name.Value, &object.Function{Parameters: node.Parameters, Body: node.Body, Scope: s}, object.FUNC)
		} else {
			return &object.Function{Parameters: node.Parameters, Body: node.Body, Scope: s}
		}
	case *ast.CallExpression:
		function := Run(node.Function, scope)
		args := runExpressions(node.Arguments, scope)

		switch function := function.(type) {
		case *object.Function:
			return applyFunction(function, args)
		case *object.Builtin:
			return function.Func(args...)
		}
	default:
		return &object.Null{}
	}

	return &object.Null{}
}

func runRebind(left *ast.Identifier, right ast.Expression, operator string, scope *object.Scope) object.Object {
	rightVal := Run(right, scope)
	leftVal, ok := scope.Get(left.Value)
	if !ok {
		panic("Undefined variable " + left.Value)
	}

	if operator == "=" {
		scope.Set(left.Value, rightVal, object.LET)
		return &object.Null{}
	}

	switch leftVal := leftVal.(type) {
	case *object.Integer:
		switch rightVal := rightVal.(type) {
		case *object.Integer:
			switch operator {
			case "+=":
				scope.Set(left.Value, &object.Integer{Value: leftVal.Value + rightVal.Value}, object.LET)
				return &object.Null{}
			case "-=":
				scope.Set(left.Value, &object.Integer{Value: leftVal.Value - rightVal.Value}, object.LET)
				return &object.Null{}
			case "*=":
				scope.Set(left.Value, &object.Integer{Value: leftVal.Value * rightVal.Value}, object.LET)
				return &object.Null{}
			case "/=":
				scope.Set(left.Value, &object.Integer{Value: leftVal.Value / rightVal.Value}, object.LET)
				return &object.Null{}
			case "%=":
				scope.Set(left.Value, &object.Integer{Value: leftVal.Value % rightVal.Value}, object.LET)
				return &object.Null{}
			}
		}
	case *object.String:
		switch rightVal := rightVal.(type) {
		case *object.String:
			switch operator {
			case "+=":
				scope.Set(left.Value, &object.String{Value: leftVal.Value + rightVal.Value}, object.LET)
				return &object.Null{}
			}
		}
	}
	return &object.Null{}
}

func runAST(statements []ast.Statement, scope *object.Scope) object.Object {
	var result object.Object
	for _, statement := range statements {
		result = Run(statement, scope)
	}

	if result == nil {
		return &object.Null{}
	}
	return result
}

func runBlockStatement(block *ast.BlockStatement, scope *object.Scope) object.Object {
	var result object.Object
	for _, statement := range block.Statements {
		result = Run(statement, scope)
	}

	if result == nil {
		return &object.Null{}
	}
	return result
}

func runExpressions(exps []ast.Expression, scope *object.Scope) []object.Object {
	var result []object.Object
	for _, e := range exps {
		evaluated := Run(e, scope)
		result = append(result, evaluated)
	}
	return result
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		extendedScope := extendFunctionScope(fn, args)
		evaluated := Run(fn.Body, extendedScope)
		return unwrapReturnValue(evaluated)
	default:
		return &object.Null{}
	}
}

func extendFunctionScope(fn *object.Function, args []object.Object) *object.Scope {
	scope := object.NewScope()
	scope.Parent = fn.Scope

	for paramIdx, param := range fn.Parameters {
		scope.Set(param.Value, args[paramIdx], object.LET)
	}

	return scope
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}

	return obj
}

func runUnop(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return runBangOp(right)
	case "-":
		return runMinusOp(right)
	default:
		return &object.Null{}
	}
}

func runBangOp(right object.Object) object.Object {
	switch right {
	case &object.Boolean{Value: true}:
		return &object.Boolean{Value: false}
	case &object.Boolean{Value: false}:
		return &object.Boolean{Value: true}
	case &object.Null{}:
		return &object.Boolean{Value: true}
	default:
		return &object.Boolean{Value: false}
	}
}

func runMinusOp(right object.Object) object.Object {
	if right.Type() != object.INTEGER {
		return &object.Null{}
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func runBinop(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER && right.Type() == object.INTEGER:
		return runIntegerBinop(operator, left, right)
	case left.Type() == object.STRING && right.Type() == object.STRING:
		return runStringBinop(operator, left, right)
	case left.Type() == object.BOOLEAN && right.Type() == object.BOOLEAN:
		return runBooleanBinop(operator, left, right)
	default:
		return &object.Null{}
	}
}

func runIntegerBinop(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value
	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "%":
		return &object.Integer{Value: leftVal % rightVal}
	case "==":
		return &object.Boolean{Value: leftVal == rightVal}
	case "!=":
		return &object.Boolean{Value: leftVal != rightVal}
	case ">":
		return &object.Boolean{Value: leftVal > rightVal}
	case "<":
		return &object.Boolean{Value: leftVal < rightVal}
	case ">=":
		return &object.Boolean{Value: leftVal >= rightVal}
	case "<=":
		return &object.Boolean{Value: leftVal <= rightVal}
	default:
		return &object.Null{}
	}
}

func runStringBinop(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	switch operator {
	case "+":
		return &object.String{Value: leftVal + rightVal}
	case "==":
		return &object.Boolean{Value: leftVal == rightVal}
	case "!=":
		return &object.Boolean{Value: leftVal != rightVal}
	default:
		return &object.Null{}
	}
}

func runBooleanBinop(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Boolean).Value
	rightVal := right.(*object.Boolean).Value

	switch operator {
	case "==":
		return &object.Boolean{Value: leftVal == rightVal}
	case "!=":
		return &object.Boolean{Value: leftVal != rightVal}
	default:
		return &object.Null{}
	}
}

func runIfExpression(ie *ast.IfExpression, scope *object.Scope) object.Object {
	condition := Run(ie.Condition, scope)
	if isTruthy(condition) {
		return Run(ie.Consequence, scope)
	} else if ie.Alternative != nil {
		return Run(ie.Alternative, scope)
	} else {
		return &object.Null{}
	}
}

func isTruthy(obj object.Object) bool {
	switch obj.Type() {
	case object.NULL:
		return false
	case object.BOOLEAN:
		switch obj.(*object.Boolean).Value {
		case true:
			return true
		case false:
			return false
		}
	}
	return false
}

func runForExpression(fe *ast.ForExpression, s *object.Scope) object.Object {
	scope := object.NewScope()
	scope.Parent = s
	var result object.Object
	Run(fe.Initializer, scope)
	for {
		condition := Run(fe.Condition, scope)
		if !isTruthy(condition) {
			break
		}

		Run(fe.Post, scope)

		result = Run(fe.Body, scope)
		if _, ok := result.(*object.ReturnValue); ok {
			break
		}
	}
	return &object.Null{}
}
