package schemer

import (
	"fmt"

	"github.com/dop251/goja"
)

type Context struct {
	vm     *goja.Runtime
	output map[string]interface{}
	ready  bool
}

func NewContext() *Context {
	ctx := &Context{
		vm:     goja.New(),
		output: make(map[string]interface{}),
	}

	ctx.initialize()

	return ctx
}

func (ctx *Context) initialize() {

	// Native functions
	console := ctx.vm.NewObject()
	console.Set("log", func(args ...interface{}) {
		fmt.Println(args...)
	})
	ctx.vm.Set("console", console)
}

func (ctx *Context) PreloadScript(p *goja.Program) error {
	_, err := ctx.vm.RunProgram(p)
	if err != nil {
		return err
	}

	ctx.ready = true

	return nil
}

func (ctx *Context) IsReady() bool {
	return ctx.ready
}
