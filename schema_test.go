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
	err := UnmarshalJSON([]byte(source), schema)
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
	assert.Equal(t, TYPE_STRING, schema.Fields["tags"].Subtype.Type)

	// Check map type
	attrs := schema.Fields["attributes"].Schema
	assert.Equal(t, TYPE_STRING, attrs.Fields["title"].Type)
	assert.Equal(t, TYPE_STRING, attrs.Fields["team"].Type)
}

func TestSchemaNormalizeWithInternalField(t *testing.T) {

	definition := `{
	"name": { "type": "string" }
}`

	data := map[string]interface{}{
		"$internal": "InternalValue",
		"name":      "Bob",
	}

	// Initializing schema
	schema := NewSchema()
	err := UnmarshalJSON([]byte(definition), schema)
	if err != nil {
		t.Error(err)
	}

	result := schema.Normalize(data)

	assert.Equal(t, "Bob", result["name"])
	assert.Equal(t, "InternalValue", result["$internal"])
}

func TestSchemaNormalizeWithNestedStructure(t *testing.T) {

	definition := `{
	"attributes": {
		"type": "map",
		"fields": {
			"title": { "type": "string" },
			"team": { "type": "string" }
		}
	},
	"nested": {
		"type": "array",
		"subtype": {
			"type": "array",
			"subtype": "string"
		}
	}
}`

	data := map[string]interface{}{
		"attributes": map[string]interface{}{
			"title": "hello",
			"team":  "world",
		},
		"nested": []interface{}{
			[]interface{}{"tag1", "tag2"},
		},
		"attributes.title": "new_hello",
		"nested.0.0":       "new_tag1",
	}

	// Initializing schema
	schema := NewSchema()
	err := UnmarshalJSON([]byte(definition), schema)
	if err != nil {
		t.Error(err)
	}

	result := schema.Normalize(data)

	// field: attributes
	attrs := result["attributes"].(map[string]interface{})
	assert.Equal(t, "hello", attrs["title"])
	assert.Equal(t, "world", attrs["team"])

	// field: nested
	nested := result["nested"].([]interface{})
	sub := nested[0].([]interface{})
	assert.Contains(t, sub, "tag1")
	assert.Contains(t, sub, "tag2")

	// fullpath field
	assert.Equal(t, "new_hello", result["attributes.title"])
	assert.Equal(t, "new_tag1", result["nested.0.0"])
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
	},
	"attachments": {
		"type": "array",
		"subtype": "map",
		"fields": {
			"filename": { "type": "string" },
			"size": { "type": "int" }
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
	},
	"attachments": [
		{ "filename": "file1.txt", "size": 123 },
		{ "filename": "file2.txt", "size": 456 }
	]
}`

	// Initializing schema
	schema := NewSchema()
	err := UnmarshalJSON([]byte(definition), schema)
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

	// field: attachments
	attachmentsValue := record.GetValue("attachments")
	assert.Equal(t, TYPE_ARRAY, attachmentsValue.Definition.Type)
	assert.Equal(t, TYPE_MAP, attachmentsValue.Definition.Subtype.Type)

	aFilenameValue := record.GetValue("attachments[0].filename")
	assert.NotNil(t, aFilenameValue)
	assert.Equal(t, TYPE_STRING, aFilenameValue.Definition.Type)
	assert.Equal(t, "file1.txt", aFilenameValue.Data)

	aSizeValue := record.GetValue("attachments[0].size")
	assert.NotNil(t, aSizeValue)
	assert.Equal(t, TYPE_INT64, aSizeValue.Definition.Type)
	assert.Equal(t, int64(123), aSizeValue.Data)
}

func TestSchema_MSSQL_Types(t *testing.T) {

	definition := `{
	"id": { "type": "uint" },
	"name": { "type": "string" },
	"age": { "type": "int" },
	"salary": { "type": "float" },
	"is_active": { "type": "bool" },
	"birth_date": { "type": "time" },
	"join_date": { "type": "time" },
	"last_updated": { "type": "time" },
	"phone_number": { "type": "string" },
	"email": { "type": "string" },
	"address": { "type": "string" },
	"notes": { "type": "string" },
	"photo": { "type": "binary" },
	"income": { "type": "float" },
	"unique_id": { "type": "string" },
	"xml_data": { "type": "string" },
	"json_data": { "type": "string" },
	"geometry_data": { "type": "binary" },
	"geography_data": { "type": "binary" },
	"created_at": { "type": "time" },
	"tinyint_col": { "type": "int" },
	"bigint_col": { "type": "int" },
	"real_col": { "type": "float" },
	"numeric_col": { "type": "float" },
	"smallmoney_col": { "type": "float" },
	"ntext_col": { "type": "string" },
	"binary_col": { "type": "binary" },
	"image_col": { "type": "binary" },
	"datetime2_col": { "type": "time" },
	"time_col": { "type": "time" },
	"timestamp_col": { "type": "time" },
	"hierarchyid_col": { "type": "any" }
}`

	recordSource := `{
	"address": "123 Main St",
	"age": 30,
	"bigint_col": 1234567890,
	"binary_col": "ASNFAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=",
	"birth_date": "1990-05-15T00:00:00Z",
	"created_at": "2023-05-05T15:45:00-07:00",
	"datetime2_col": "2023-05-05T12:00:00Z",
	"email": "john@example.com",
	"geography_data": "5hAAAAEMqvHSTWKAUsBeS8gHPVtEQA==",
	"geometry_data": "AAAAAAEMAAAAAAAA8D8AAAAAAADwPw==",
	"hierarchyid_col": "WsA=",
	"id": 2,
	"image_col": "/+7dzA==",
	"income": 5000.75,
	"is_active": true,
	"join_date": "2022-01-01T09:00:00Z",
	"json_data": "{\"key\": \"value\"}",
	"last_updated": "2023-05-05T12:30:00Z",
	"name": "John Doe",
	"notes": "Some notes",
	"ntext_col": "Some text",
	"numeric_col": "MTIzNDUuNjc=",
	"phone_number": "1234567890",
	"photo": "ASNFZ4mrze8=",
	"real_col": 3.140000104904175,
	"salary": "MjUwMC41MA==",
	"smallmoney_col": "MTAwLjI1MDA=",
	"time_col": "0001-01-01T15:30:00Z",
	"timestamp_col": "AAAAAAAAB9E=",
	"tinyint_col": 5,
	"unique_id": "/xmWb4aLEdC0LQDAT8lk/w==",
	"xml_data": "<root><data>Some XML data</data></root>"
}`

	// Initializing schema
	schema := NewSchema()
	err := UnmarshalJSON([]byte(definition), schema)
	if err != nil {
		t.Error(err)
	}

	// Initializing record
	var rawData map[string]interface{}
	json.Unmarshal([]byte(recordSource), &rawData)

	// Scan raw data
	_ = schema.Scan(rawData)

}
