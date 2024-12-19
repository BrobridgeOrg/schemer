package schemer

import (
	"testing"
)

func BenchmarkSchemaNormalize(b *testing.B) {

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
		schema.Normalize(rawData)
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
