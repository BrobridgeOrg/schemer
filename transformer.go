package schemer

import (
	"sync"
)

type Transformer struct {
	source  *Schema
	dest    *Schema
	script  string
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

	return t
}

func (t *Transformer) Transform(env map[string]interface{}, input map[string]interface{}) ([]map[string]interface{}, error) {

	data := t.source.Normalize(input)

	// Preparing context and runtime
	ctx := t.ctxPool.Get().(*Context)
	defer t.ctxPool.Put(ctx)
	if !ctx.IsReady() {
		err := ctx.PreloadScript(t.script)
		if err != nil {
			return nil, err
		}
	}

	ctx.vm.Set("env", env)
	ctx.vm.Set("source", data)

	//var fn func() map[string]interface{}
	var fn func() interface{}
	err := ctx.vm.ExportTo(ctx.vm.Get("main"), &fn)
	if err != nil {
		return nil, err
	}

	result := fn()
	if result == nil {
		return nil, nil
	}

	// Result is an object
	if v, ok := result.(map[string]interface{}); ok {

		// Normalized for destination schema then returning result
		return []map[string]interface{}{
			t.dest.Normalize(v),
		}, nil
	} else if v, ok := result.([]interface{}); ok {
		// Result is an array

		returnedValue := make([]map[string]interface{}, len(v))
		for i, d := range v {

			if v, ok := d.(map[string]interface{}); ok {
				returnedValue[i] = v
			}
		}

		return returnedValue, nil
	}

	return nil, nil
}

func (t *Transformer) SetScript(script string) {
	t.script = `function main() {` + script + `}`
}
