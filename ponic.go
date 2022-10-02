package main

import (
	"bufio"
	"os"

	"github.com/danecwalker/ponic/engine/lexer"
	"github.com/danecwalker/ponic/engine/object"
	"github.com/danecwalker/ponic/engine/parser"
	"github.com/danecwalker/ponic/engine/runtime"
)

func main() {
	f, err := os.Open("./examples/example.pc")
	if err != nil {
		panic(err)
	}

	l := lexer.NewLexer(bufio.NewReader(f))
	p := parser.NewParser(l)

	global_scope := object.NewScope()
	runtime.Run(p.Parse(), global_scope)
}
