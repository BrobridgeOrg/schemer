package v8go_runtime

import (
	"encoding/base64"
	"fmt"
	"sync/atomic"

	_ "embed"

	"github.com/BrobridgeOrg/schemer"
	jsoniter "github.com/json-iterator/go"
	"rogchap.com/v8go"
)

//go:embed js/msgpack.js
var msgpackJS string

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Runtime struct {
	counter int64
	iso     *v8go.Isolate
	ctx     *v8go.Context
	mainFn  *v8go.Function
}

func NewRuntime() *Runtime {

	r := &Runtime{}

	r.iso = v8go.NewIsolate()
	r.ctx = v8go.NewContext(r.iso)

	// add console.log native function
	console := v8go.NewObjectTemplate(r.iso)
	console.Set("log", v8go.NewFunctionTemplate(r.iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		/*
			/*
				args := make([]interface{}, len(info.Args()))
				for i, arg := range info.Args() {
					args[i] = arg.String()
				}
			args := info.Args()
		*/

		args := make([]interface{}, len(info.Args()))
		for i, arg := range info.Args() {
			args[i] = arg.Object()
		}

		fmt.Println(args...)

		return nil
	}))

	consoleObj, _ := console.NewInstance(r.ctx)
	r.ctx.Global().Set("console", consoleObj)

	// Add built-infunctions
	native := v8go.NewObjectTemplate(r.iso)
	native.Set("toBase64", v8go.NewFunctionTemplate(r.iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {

		source := info.Args()[0].Object()

		length, err := source.Get("length")
		if err != nil {
			panic(err)
		}

		data := make([]byte, length.Int32())
		for i := 0; i < int(length.Int32()); i++ {

			b, err := source.GetIdx(uint32(i))
			if err != nil {
				panic(err)
			}

			data[i] = byte(b.Int32())
		}

		encodedData := base64.StdEncoding.EncodeToString(data)

		v, err := v8go.NewValue(r.iso, encodedData)
		if err != nil {
			panic(err)
		}

		return v
	}))

	nativeObj, _ := native.NewInstance(r.ctx)
	r.ctx.Global().Set("native", nativeObj)

	r.LoadScript(msgpackJS)

	//defer iso.Dispose()
	return r
}

func (r *Runtime) SetEnv(value map[string]interface{}) {

	if value == nil {
		env := v8go.NewObjectTemplate(r.iso)
		envObj, _ := env.NewInstance(r.ctx)
		r.ctx.Global().Set("env", envObj)
		return
	}

	obj := r.convertToV8Object(value)
	r.ctx.Global().Set("env", obj)
}

func (r *Runtime) LoadScript(script string) error {

	counter := atomic.AddInt64(&r.counter, 1)

	_, err := r.ctx.RunScript(script, fmt.Sprintf("load_%d.js", counter))
	if err != nil {
		return err
	}

	return nil
}

func (r *Runtime) Compile(script string) error {
	return r.LoadScript(script)
}

func (r *Runtime) Execute(sourceSchema *schemer.Schema, data map[string]interface{}) ([]map[string]interface{}, error) {
	/*
		// Normalize data for JavaScript
		if sourceSchema != nil {
			r.normalize(sourceSchema, data)
		}
	*/

	//	fmt.Println(data["binary"].([]interface{})[0])

	// Convert data to V8 object
	value := r.convertToV8Object(data)
	/*
			// Convert result to map
			_, err := r.convertV8ObjectToStruct(obj)
			if err != nil {
				panic(err)
			}
		ov, _ := obj.Object().Get("binary")
		fmt.Println("V8 OBJ", ov.String())
	*/

	// Get main function from VM
	if r.mainFn == nil {
		fn, err := r.ctx.Global().Get("main")
		if err != nil {
			return nil, err
		}

		main, err := fn.AsFunction()
		if err != nil {
			return nil, err
		}

		r.mainFn = main
	}

	resultValue, err := r.mainFn.Call(r.ctx.Global(), value)
	if err != nil {
		return nil, err
	}

	if resultValue.IsNullOrUndefined() {
		return nil, nil
	}

	// Convert result to map
	result, err := r.convertV8ObjectToStruct(resultValue)
	if err != nil {
		return nil, err
	}

	// returned data is an array
	switch d := result.(type) {
	case []interface{}:

		// Prepare array
		returnedValues := make([]map[string]interface{}, len(d))

		for i, ele := range d {
			v := ele.(map[string]interface{})

			returnedValues[i] = v
		}

		return returnedValues, nil
	case map[string]interface{}:

		return []map[string]interface{}{
			d,
		}, nil

	}

	return []map[string]interface{}{}, nil
}
