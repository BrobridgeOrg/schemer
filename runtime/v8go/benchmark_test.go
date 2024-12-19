package v8go_runtime

import (
	"testing"

	"github.com/BrobridgeOrg/schemer"
	"rogchap.com/v8go"
)

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

	iso := v8go.NewIsolate()
	ctx := v8go.NewContext(iso)

	ctx.RunScript("function main() {}", "main.js")
	fn, _ := ctx.Global().Get("main")
	main, _ := fn.AsFunction()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		main.Call(ctx.Global(), v8go.Undefined(iso))
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
