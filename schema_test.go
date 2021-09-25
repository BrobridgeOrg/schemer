package schemer

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSchemaUnmarshal(t *testing.T) {

	source := `{
	"name": { "type": "string" },
	"image": { "type": "binary" },
	"enabled": { "type": "bool" },
	"age": { "type": "uint" },
	"balance": { "type": "int" },
	"score": { "type": "float" },
	"createdAt": { "type": "time" },
	"tags": {
		"type": "array",
		"subtype": "string"
	},
	"attributes": {
		"type": "map",
		"fields": {
			"title": { "type": "string" },
			"team": { "type": "string" }
		}
	},
	"attached": { "type": "any" }
}`

	schema := NewSchema()
	err := UnmarsalJSON([]byte(source), schema)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, TYPE_STRING, schema.Fields["name"].Type)
	assert.Equal(t, TYPE_BINARY, schema.Fields["image"].Type)
	assert.Equal(t, TYPE_BOOLEAN, schema.Fields["enabled"].Type)
	assert.Equal(t, TYPE_UINT64, schema.Fields["age"].Type)
	assert.Equal(t, TYPE_INT64, schema.Fields["balance"].Type)
	assert.Equal(t, TYPE_FLOAT64, schema.Fields["score"].Type)
	assert.Equal(t, TYPE_TIME, schema.Fields["createdAt"].Type)
	assert.Equal(t, TYPE_ARRAY, schema.Fields["tags"].Type)
	assert.Equal(t, TYPE_MAP, schema.Fields["attributes"].Type)
	assert.Equal(t, TYPE_ANY, schema.Fields["attached"].Type)

	// Check array type
	assert.Equal(t, TYPE_STRING, schema.Fields["tags"].Subtype)

	// Check map type
	attrs := schema.Fields["attributes"].Definition
	assert.Equal(t, TYPE_STRING, attrs.Fields["title"].Type)
	assert.Equal(t, TYPE_STRING, attrs.Fields["team"].Type)
}

func TestSchemaScan(t *testing.T) {

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
	err := UnmarsalJSON([]byte(definition), schema)
	if err != nil {
		t.Error(err)
	}

	// Initializing record
	var rawData map[string]interface{}
	json.Unmarshal([]byte(recordSource), &rawData)

	// Scan raw data
	record := schema.Scan(rawData)

	// field: name
	value := record.GetValue("name")
	assert.Equal(t, TYPE_STRING, value.Definition.Type)
	assert.Equal(t, "Fred", value.Data.(string))

	// field: balance
	balanceValue := record.GetValue("balance")
	assert.Equal(t, TYPE_INT64, balanceValue.Definition.Type)
	assert.Equal(t, int64(123456), balanceValue.Data.(int64))

	// field: key
	keyValue := record.GetValue("key")
	assert.Equal(t, TYPE_BINARY, keyValue.Definition.Type)
	if bytes.Compare([]byte{byte(12), byte(34)}, keyValue.Data.([]byte)) != 0 {
		t.Log(keyValue.Data)
		t.Error("key is not equal to expected value")
	}

	// field: createdAt
	createdAtValue := record.GetValue("createdAt")
	assert.Equal(t, TYPE_TIME, createdAtValue.Definition.Type)
	assert.Equal(t, int64(1595182568), createdAtValue.Data.(time.Time).Unix())

	// field: updatedAt
	updatedAtValue := record.GetValue("updatedAt")
	assert.Equal(t, TYPE_TIME, updatedAtValue.Definition.Type)
	assert.Equal(t, int64(1595182568), updatedAtValue.Data.(time.Time).Unix())

	// field: attributes.title
	titleValue := record.GetValue("attributes.title")
	assert.Equal(t, TYPE_STRING, titleValue.Definition.Type)
	assert.Equal(t, "Architect", titleValue.Data.(string))

	// field: attributes.team
	teamValue := record.GetValue("attributes.team")
	assert.Equal(t, TYPE_STRING, teamValue.Definition.Type)
	assert.Equal(t, "product", teamValue.Data.(string))
}
