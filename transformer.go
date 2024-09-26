package schemer

import (
	"reflect"
	"sync"
	"time"

	"github.com/BrobridgeOrg/schemer/types"
	"github.com/dop251/goja"
)

type Transformer struct {
	source  *Schema
	dest    *Schema
	script  string
	program *goja.Program
	ctxPool sync.Pool
}

func NewTransformer(source *Schema, dest *Schema) *Transformer {

	t := &Transformer{
		source: source,
		dest:   dest,
		script: `function main() { return source; }`,
	}

	t.ctxPool.New = func() interface{} {
		return NewContext()
	}

	t.program, _ = goja.Compile("transformer", t.script, false)

	return t
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

			if def.Info.(*types.Time).Percision != types.TIME_PERCISION_MICROSECOND {
				v, _ := ctx.vm.New(ctx.vm.Get("Date").ToObject(ctx.vm), ctx.vm.ToValue(val.(time.Time).UnixMicro()/1e3))
				data[fieldName] = v
			}
			continue
		}
	}
}

func (t *Transformer) initializeContext(ctx *Context, env map[string]interface{}, schema *Schema, data map[string]interface{}) error {

	if !ctx.IsReady() {
		err := ctx.PreloadScript(t.program)
		if err != nil {
			return err
		}
	}

	// Initializing environment varable
	ctx.vm.Set("env", env)

	// Normorlize data for JavaScript
	if t.source != nil {
		t.normalize(ctx, t.source, data)
	}

	ctx.vm.Set("source", data)

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

func (t *Transformer) runScript(ctx *Context) ([]map[string]interface{}, error) {

	var fn func() interface{}
	err := ctx.vm.ExportTo(ctx.vm.Get("main"), &fn)
	if err != nil {
		return nil, err
	}

	result := fn()
	if result == nil {
		return nil, nil
	}

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

	// Preparing context and runtime
	ctx := t.ctxPool.Get().(*Context)
	defer t.ctxPool.Put(ctx)

	err := t.initializeContext(ctx, env, t.source, data)
	if err != nil {
		return nil, err
	}

	// Run script to process data
	result, err := t.runScript(ctx)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, nil
	}

	return result, nil
}

func (t *Transformer) SetScript(script string) error {
	t.script = `
function run() {` + script + `}
function scanStruct(obj) {
	for (key in obj) {
		val = obj[key]
		if (val === undefined) {
			delete obj[key]
		} else if (val == null) {
			continue
		} else if (val instanceof Array) {
			scanStruct(val)
		} else if (typeof val === 'object') {
			scanStruct(val)
		}
	}
}
function main() {
	v = run()
	if (v === null)
		return null
	scanStruct(v)
	return v
}
`

	p, err := goja.Compile("transformer", t.script, false)
	if err != nil {
		return err
	}

	t.program = p

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
