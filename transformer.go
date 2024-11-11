package schemer

import (
	_ "embed"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/BrobridgeOrg/schemer/types"
	"github.com/dop251/goja"
)

//go:embed js/dummy.js
var dummyJS string

//go:embed js/core.js
var coreJS string

type Transformer struct {
	source        *Schema
	dest          *Schema
	script        string
	program       *goja.Program
	ctx           *Context
	relationships map[string][]string
}

func NewTransformer(source *Schema, dest *Schema) *Transformer {

	t := &Transformer{
		source:        source,
		dest:          dest,
		relationships: make(map[string][]string),
	}

	// Preparing context
	ctx, err := t.createContext()
	if err != nil {
		panic(err)
	}

	t.ctx = ctx

	t.injectFuncs()

	err = t.SetScript(`return source`)
	if err != nil {
		panic(err)
	}

	return t
}

func (t *Transformer) createContext() (*Context, error) {

	ctx := NewContext()

	err := ctx.LoadScript(dummyJS)
	if err != nil {
		return nil, err
	}

	err = ctx.LoadScript(coreJS)
	if err != nil {
		return nil, err
	}

	return ctx, nil
}

func (t *Transformer) normalize(ctx *Context, schema *Schema, data map[string]interface{}) {

	for fieldName, def := range schema.Fields {

		val, ok := data[fieldName]
		if !ok {
			continue
		}

		if def.Type == TYPE_MAP {
			t.normalize(ctx, def.Schema, val.(map[string]interface{}))
			continue
		}

		if def.Type == TYPE_TIME {

			// Skip null
			if val == nil {
				continue
			}

			if def.Info.(*types.Time).Precision != types.TIME_PRECISION_MICROSECOND {
				v, _ := ctx.vm.New(ctx.vm.Get("Date").ToObject(ctx.vm), ctx.vm.ToValue(val.(time.Time).UnixMicro()/1e3))
				data[fieldName] = v
			}
			continue
		}
	}
}

func (t *Transformer) initializeContext(ctx *Context, env map[string]interface{}, schema *Schema, data map[string]interface{}) error {

	if !ctx.IsReady() {
		return fmt.Errorf("Context is not ready")
	}

	// Initializing environment varable
	ctx.vm.Set("env", env)

	// Normorlize data for JavaScript
	if t.source != nil {
		t.normalize(ctx, t.source, data)
	}

	//	ctx.vm.Set("source", data)

	return nil
}

func (t *Transformer) handleArrayValue(arrayValue []interface{}) error {

	for i, value := range arrayValue {
		v := reflect.ValueOf(value)
		switch v.Kind() {
		/*
			case reflect.Slice:
				err := t.handleArrayValue(value.([]interface{}))
				if err != nil {
					return err
				}
		*/
		case reflect.Map:
			err := t.handleMapValue(value.(map[string]interface{}))
			if err != nil {
				return err
			}
		default:
			// Convert Data object to time.Time
			switch d := value.(type) {
			case *goja.Object:
				if value.(*goja.Object).ClassName() == "Date" {
					arrayValue[i] = d.Export()
				}
			}
		}
	}

	return nil
}

func (t *Transformer) handleMapValue(returnedValue map[string]interface{}) error {

	for key, value := range returnedValue {
		v := reflect.ValueOf(value)
		switch v.Kind() {
		/*
			case reflect.Slice:
				err := t.handleArrayValue(value.([]interface{}))
				if err != nil {
					return err
				}
		*/
		case reflect.Map:
			err := t.handleMapValue(value.(map[string]interface{}))
			if err != nil {
				return err
			}
		default:
			// Convert Data object to time.Time
			switch d := value.(type) {
			case *goja.Object:
				if value.(*goja.Object).ClassName() == "Date" {
					returnedValue[key] = d.Export()
				}
			}
		}
	}

	return nil
}

func (t *Transformer) normalizeValue(v map[string]interface{}) (map[string]interface{}, error) {

	var val map[string]interface{}

	// Normalize for destination schema if it exists
	if t.dest != nil {
		val = t.dest.Normalize(v)
	} else if t.source != nil {

		// Inherit source schema
		val = t.source.Normalize(v)
	} else {
		val = v
	}

	return val, nil
}

func (t *Transformer) prepareRefs(source map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for sourceKey, value := range source {
		keyParts := strings.Split(sourceKey, ".")
		level := result

		for i := 0; i < len(keyParts); i++ {
			part := keyParts[i]

			// If we are at the last part, assign the value
			if i == len(keyParts)-1 {
				level[part] = value
			} else {
				// If the part does not exist, create a new map
				if _, ok := level[part]; !ok {
					level[part] = make(map[string]interface{})
				}

				// Move deeper into the nested map
				nextLevel, _ := level[part].(map[string]interface{})
				level = nextLevel
			}
		}
	}

	return result
}

func (t *Transformer) injectFuncs() error {

	t.ctx.vm.Set("normalize", func(call goja.FunctionCall) goja.Value {
		input := call.Argument(0).Export().(map[string]interface{})
		result := t.prepareRefs(input)
		return t.ctx.vm.ToValue(result)
	})

	return nil
}

func (t *Transformer) runScript(ctx *Context, data map[string]interface{}) ([]map[string]interface{}, error) {

	main, ok := goja.AssertFunction(ctx.vm.Get("main"))
	if !ok {
		return nil, fmt.Errorf("main is not a function")
	}

	source := ctx.vm.ToValue(data)
	res, err := main(goja.Undefined(), source)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	/*
		var fn func() interface{}
		err := ctx.vm.ExportTo(ctx.vm.Get("main"), &fn)
		if err != nil {
			return nil, err
		}

		result := fn()
	*/

	var result interface{}
	if goja.IsNull(res) || goja.IsUndefined(res) {
		return nil, nil
	}

	result = res.Export()

	v := reflect.ValueOf(result)
	switch v.Kind() {
	case reflect.Slice:

		v := result.([]interface{})

		// Prepare array
		returnedValues := make([]map[string]interface{}, len(v))

		for i, d := range v {
			v := d.(map[string]interface{})

			// Deal with JavaScript Object
			err := t.handleMapValue(v)
			if err != nil {
				return nil, err
			}

			// Normalize returned data based on schema
			val, err := t.normalizeValue(v)
			if err != nil {
				return nil, err
			}

			returnedValues[i] = val
		}

		return returnedValues, nil

	default:

		v := result.(map[string]interface{})

		// Deal with JavaScript Object
		err := t.handleMapValue(v)
		if err != nil {
			return nil, err
		}

		// Normalize returned data based on schema
		val, err := t.normalizeValue(v)
		if err != nil {
			return nil, err
		}

		return []map[string]interface{}{
			val,
		}, nil
	}
}

func (t *Transformer) Transform(env map[string]interface{}, input map[string]interface{}) ([]map[string]interface{}, error) {

	var data map[string]interface{} = input
	if t.source != nil {
		data = t.source.Normalize(input)
	}

	err := t.initializeContext(t.ctx, env, t.source, data)
	if err != nil {
		return nil, err
	}

	// Run script to process data
	result, err := t.runScript(t.ctx, data)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, nil
	}

	return result, nil
}

func (t *Transformer) prepareScript(script string) string {
	return `function script(source) {` + script + `}`
}

func (t *Transformer) SetScript(script string) error {

	t.script = t.prepareScript(script)

	p, err := goja.Compile("transformer", t.script, false)
	if err != nil {
		return err
	}

	t.program = p

	err = t.ctx.PreloadScript(t.program)
	if err != nil {
		panic(err)
	}

	return nil
}

func (t *Transformer) SetSourceSchema(schema *Schema) {
	t.source = schema
}

func (t *Transformer) GetSourceSchema() *Schema {
	return t.source
}

func (t *Transformer) SetDestinationSchema(schema *Schema) {
	t.dest = schema
}

func (t *Transformer) GetDestinationSchema() *Schema {

	if t.source != nil && t.dest == nil {
		return t.source
	}

	return t.dest
}
