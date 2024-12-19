package goja_runtime

import (
	"fmt"

	"github.com/BrobridgeOrg/schemer"
	"github.com/dop251/goja"
)

type Runtime struct {
	vm   *goja.Runtime
	main goja.Callable
}

func NewRuntime() *Runtime {

	r := &Runtime{
		vm:   goja.New(),
		main: nil,
	}

	r.initBuiltInFunctions()

	return r
}

func (r *Runtime) SetEnv(value map[string]interface{}) {
	r.vm.Set("env", value)
}

func (r *Runtime) LoadScript(script string) error {
	_, err := r.vm.RunString(script)
	if err != nil {
		return err
	}

	return nil
}

func (r *Runtime) Compile(script string) error {

	err := r.initNativeFunctions()
	if err != nil {
		return err
	}

	p, err := goja.Compile("runtime", script, false)
	if err != nil {
		return err
	}

	_, err = r.vm.RunProgram(p)
	if err != nil {
		return err
	}

	return nil
}

func (r *Runtime) Execute(sourceSchema *schemer.Schema, data map[string]interface{}) ([]map[string]interface{}, error) {

	// Normalize data for JavaScript
	if sourceSchema != nil {
		r.normalize(sourceSchema, data)
	}

	if r.main == nil {
		// Get main function from VM
		main, ok := goja.AssertFunction(r.vm.Get("main"))
		if !ok {
			return nil, fmt.Errorf("main is not a function")
		}

		r.main = main
	}

	/*
		// Preparing $ref
		ref := t.internalImpl.PrepareRefs(data)
		data["$ref"] = ref
	*/
	// Execute
	res, err := r.main(goja.Undefined(), r.vm.ToValue(data))
	if err != nil {
		return nil, err
	}

	if goja.IsNull(res) || goja.IsUndefined(res) || goja.IsNaN(res) || goja.IsInfinity(res) {
		return nil, nil
	}

	var result interface{} = res.Export()

	// returned data is an array
	switch d := result.(type) {
	case []interface{}:

		// Prepare array
		returnedValues := make([]map[string]interface{}, len(d))

		for i, ele := range d {
			v := ele.(map[string]interface{})

			// Deal with JavaScript Object
			err := r.handleMapValue(v)
			if err != nil {
				return nil, err
			}

			returnedValues[i] = v
		}

		return returnedValues, nil
	case map[string]interface{}:

		// Deal with JavaScript Object
		err = r.handleMapValue(d)
		if err != nil {
			return nil, err
		}

		return []map[string]interface{}{
			d,
		}, nil

	}

	return []map[string]interface{}{}, nil

	/*
	   if reflect.ValueOf(result).Kind() == reflect.Slice {

	   		v := result.([]interface{})

	   		// Prepare array
	   		returnedValues := make([]map[string]interface{}, len(v))

	   		for i, d := range v {
	   			v := d.(map[string]interface{})

	   			// Deal with JavaScript Object
	   			err := r.handleMapValue(v)
	   			if err != nil {
	   				return nil, err
	   			}

	   			returnedValues[i] = v
	   		}

	   		return returnedValues, nil
	   	}

	   // returned data is an object
	   v := result.(map[string]interface{})

	   // Deal with JavaScript Object
	   err = r.handleMapValue(v)

	   	if err != nil {
	   		return nil, err
	   	}

	   	return []map[string]interface{}{
	   		v,
	   	}, nil
	*/
}
