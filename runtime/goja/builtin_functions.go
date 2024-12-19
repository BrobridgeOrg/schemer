package goja_runtime

import "fmt"

func (r *Runtime) initBuiltInFunctions() {

	// Native functions
	console := r.vm.NewObject()
	console.Set("log", func(args ...interface{}) {
		fmt.Println(args...)
	})

	r.vm.Set("console", console)
}
