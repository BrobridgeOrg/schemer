package schemer

import (
	_ "embed"
	"strings"
)

//go:embed js/dummy.js
var dummyJS string

//go:embed js/core.js
var coreJS string

type TransformerOpt func(*Transformer)

func WithRuntime(runtime Runtime) func(*Transformer) {
	return func(t *Transformer) {
		t.runtime = runtime
	}
}

type Transformer struct {
	source      *Schema
	dest        *Schema
	runtime     Runtime
	passThrough bool
}

func NewTransformer(source *Schema, dest *Schema, opts ...TransformerOpt) *Transformer {

	t := &Transformer{
		source:      source,
		dest:        dest,
		runtime:     nil,
		passThrough: true,
	}

	for _, opt := range opts {
		opt(t)
	}

	// Preload scripts
	t.runtime.LoadScript(dummyJS)
	t.runtime.LoadScript(coreJS)

	return t
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

func (t *Transformer) runScript(data map[string]interface{}) ([]map[string]interface{}, error) {

	results, err := t.runtime.Execute(t.source, data)
	if err != nil {
		return nil, err
	}

	for i, result := range results {

		// Normalize returned data based on schema
		val, err := t.normalizeValue(result)
		if err != nil {
			return nil, err
		}

		results[i] = val
	}

	return results, err
}

func (t *Transformer) Reset() {
	t.passThrough = true
}

func (t *Transformer) Transform(env map[string]interface{}, input map[string]interface{}) ([]map[string]interface{}, error) {
	// Pass through if no script is set
	if t.passThrough {
		data, err := t.normalizeValue(input)
		if err != nil {
			return nil, err
		}

		return []map[string]interface{}{data}, nil
	}

	var data map[string]interface{} = input
	if t.source != nil {
		data = t.source.Normalize(input)
	}

	t.runtime.SetEnv(env)

	// Run script to process data
	result, err := t.runScript(data)
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

	// Pass through if no script is set
	SCRIPT := strings.Trim(script, " ")
	if SCRIPT == "return source" || SCRIPT == "return source;" {
		t.passThrough = true
		return nil
	}

	fullScript := t.prepareScript(script)

	err := t.runtime.Compile(fullScript)
	if err != nil {
		return err
	}

	t.passThrough = false

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
