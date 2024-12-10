package goja_runtime

import (
	"testing"

	"github.com/BrobridgeOrg/schemer"
	"github.com/dop251/goja"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

var testSchema = `{
	"string": { "type": "string" },
	"binary": { "type": "binary" },
	"int": { "type": "int" },
	"uint": { "type": "uint" },
	"float": { "type": "float" },
	"bool": { "type": "bool" },
	"time": { "type": "time" },
	"microTime": {
		"type": "time",
		"precision": "microsecond"
	},
	"object": {
		"type": "map",
		"fields": {
			"title": { "type": "string" },
			"team": { "type": "string" }
		}
	}
}`

func BenchmarkJavaScriptVM(b *testing.B) {

	vm := goja.New()

	p, _ := goja.Compile("transformer", "function main() {}", false)
	vm.RunProgram(p)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		main, ok := goja.AssertFunction(vm.Get("main"))
		if !ok {
			panic("main is not a function")
		}

		res, err := main(goja.Undefined())
		if err != nil {
			panic(err)
		}

		res.Export()
	}
}

func BenchmarkTransformer(b *testing.B) {

	testSourceSchema := schemer.NewSchema()
	err := schemer.UnmarshalJSON([]byte(testSchema), testSourceSchema)
	if err != nil {
		b.Fatal(err)
	}

	testDestSchema := schemer.NewSchema()
	err = schemer.UnmarshalJSON([]byte(testSchema), testDestSchema)
	if err != nil {
		b.Fatal(err)
	}

	// Create Runtime
	r := NewRuntime()

	// Create transformer
	transformer := schemer.NewTransformer(
		testSourceSchema,
		testDestSchema,
		schemer.WithRuntime(r),
	)

	// Set transform script
	err = transformer.SetScript(`
	return {
		string: source.string + 'TEST',
		int: source.int + 1,
		uint: source.uint + 1,
		float: source.float,
		bool: source.bool
	}
`)
	if err != nil {
		b.Fatal(err)
	}

	// Transform
	rawData := `{
	"string": "Brobridge",
	"int": -9527,
	"uint": 9527,
	"float": 11.15,
	"bool": false
}`
	var sourceData map[string]interface{}
	err = json.Unmarshal([]byte(rawData), &sourceData)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := transformer.Transform(nil, sourceData)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTransformer_PassThrough(b *testing.B) {

	testSourceSchema := schemer.NewSchema()
	err := schemer.UnmarshalJSON([]byte(testSchema), testSourceSchema)
	if err != nil {
		b.Fatal(err)
	}

	testDestSchema := schemer.NewSchema()
	err = schemer.UnmarshalJSON([]byte(testSchema), testDestSchema)
	if err != nil {
		b.Fatal(err)
	}

	// Create Runtime
	r := NewRuntime()

	// Create transformer
	transformer := schemer.NewTransformer(
		testSourceSchema,
		testDestSchema,
		schemer.WithRuntime(r),
	)

	// Set transform script
	err = transformer.SetScript(`return source`)
	if err != nil {
		b.Fatal(err)
	}

	// Transform
	rawData := `{
	"string": "Brobridge",
	"int": -9527,
	"uint": 9527,
	"float": 11.15,
	"bool": false
}`
	var sourceData map[string]interface{}
	err = json.Unmarshal([]byte(rawData), &sourceData)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := transformer.Transform(nil, sourceData)
		if err != nil {
			b.Fatal(err)
		}
	}
}
