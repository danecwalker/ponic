// Copyright (c) 2022 DevDane <dane@danecwalker.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package runtime

import (
	"fmt"

	"github.com/danecwalker/ponic/engine/object"
)

type Builtin = func(args ...object.Object) object.Object

var Builtins = map[string]Builtin{
	"print": func(args ...object.Object) object.Object {
		_args := make([]interface{}, len(args))
		for i, arg := range args {
			_args[i] = unescape(arg.Inspect())
		}
		fmt.Print(_args...)
		return nil
	},
}
