// Copyright (c) 2022 DevDane <dane@danecwalker.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package runtime

import (
	"fmt"
	"strconv"

	"github.com/danecwalker/ponic/engine/object"
)

type Builtin = func(args ...object.Object) object.Object

var Builtins = map[string]Builtin{
	"print": _print,

	"scan": _scan,

	"int": _int,
}

func _print(args ...object.Object) object.Object {
	_args := make([]interface{}, len(args))
	for i, arg := range args {
		_args[i] = unescape(arg.Inspect())
	}
	fmt.Print(_args...)
	return nil
}

func _scan(args ...object.Object) object.Object {
	var input string
	_print(args...)
	fmt.Scan(&input)
	return &object.String{Value: input}
}

func _int(args ...object.Object) object.Object {
	if len(args) != 1 {
		panic("wrong number of arguments.")
	}

	switch arg := args[0].(type) {
	case *object.Integer:
		return arg
	case *object.String:
		i, err := strconv.ParseInt(arg.Value, 10, 64)
		if err != nil {
			panic(err)
		}
		return &object.Integer{Value: i}
	default:
		panic("argument to `int` not supported.")
	}
}
