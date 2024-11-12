package schemer

import (
	"testing"

	"github.com/dop251/goja"
)

func BenchmarkJavaScriptVM(b *testing.B) {

	vm := goja.New()

	p, _ := goja.Compile("transformer", "function main() {}", false)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		vm.RunProgram(p)

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

func BenchmarkScan(b *testing.B) {

	definition := `{
	"name": { "type": "string" },
	"balance": { "type": "int" },
	"key": { "type": "binary" },
	"createdAt": { "type": "time" },
	"updatedAt": { "type": "time" },
	"attributes": {
		"type": "map",
		"fields": {
			"title": { "type": "string" },
			"team": { "type": "string" }
		}
	}
}`

	recordSource := `{
	"name": "Fred",
	"balance": 123456,
	"key": [ 12, 34 ],
	"createdAt": 1595182568,
	"updatedAt": "2020-07-19T18:16:08.000001Z",
	"attributes": {
		"title": "Architect",
		"team": "product"
	}
}`

	// Initializing schema
	schema := NewSchema()
	err := UnmarshalJSON([]byte(definition), schema)
	if err != nil {
		b.Fatal(err)
	}

	// Initializing record
	var rawData map[string]interface{}
	json.Unmarshal([]byte(recordSource), &rawData)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Scan raw data
		schema.Scan(rawData)
	}
}

func BenchmarkTransformer(b *testing.B) {

	testSourceSchema := NewSchema()
	err := UnmarshalJSON([]byte(testSource), testSourceSchema)
	if err != nil {
		b.Fatal(err)
	}

	testDestSchema := NewSchema()
	err = UnmarshalJSON([]byte(testDest), testDestSchema)
	if err != nil {
		b.Fatal(err)
	}

	// Create transformer
	transformer := NewTransformer(testSourceSchema, testDestSchema)

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
