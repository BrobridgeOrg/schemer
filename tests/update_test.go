package tests

import (
	"testing"

	"github.com/BrobridgeOrg/schemer"
	"github.com/stretchr/testify/assert"
)

var testSource = `{
	"string": { "type": "string" },
	"any": { "type": "any" },
	"microTime": {
		"type": "time",
		"precision": "microsecond"
	},
	"object": {
		"type": "map",
		"fields": {
			"title": { "type": "string" },
			"team": { "type": "string" },
			"tags": {
				"type": "array",
				"subtype": "string"
			},
			"multidimensionalArray": {
				"type": "array",
				"subtype": {
					"type": "array",
					"subtype": "string"
				}
			},
			"authors": {
				"type": "array",
				"subtype": "map",
				"fields": {
					"name": { "type": "string" },
					"email": { "type": "string" }
				}
			},
			"nestedObject": {
				"type": "map",
				"fields": {
					"title": { "type": "string" },
					"desc": { "type": "string" }
				}
			}
		}
	}
}`

func Test_NestedObject_Updates(t *testing.T) {

	testSourceSchema := schemer.NewSchema()
	err := schemer.UnmarshalJSON([]byte(testSource), testSourceSchema)
	if err != nil {
		t.Error(err)
	}

	source := testSourceSchema.Normalize(map[string]interface{}{
		"string":       "Test String",
		"object.id":    "NotExists",
		"object.title": "Test Title",
		"object.team":  "A Team",
		"object.tags": []string{
			"tag1",
			"tag2",
		},
		"object.nestedObject.title": "Test Title",
		"object.authors": []map[string]interface{}{
			{
				"name":  "John Doe",
				"email": "test@example.com",
			},
		},
		"object.multidimensionalArray": [][]string{
			{"1", "2", "3"},
			{"4", "5", "6"},
		},
		"any": "Any Value",
	})

	// Normal fields
	assert.Nil(t, source["object.id"])
	assert.Equal(t, "Test String", source["string"].(string))
	assert.Equal(t, "Any Value", source["any"].(string))

	// properties of object
	assert.Equal(t, "Test Title", source["object.title"].(string))
	assert.Equal(t, "A Team", source["object.team"].(string))

	// Array of strings
	assert.Len(t, source["object.tags"], 2)
	tags := source["object.tags"].([]interface{})
	assert.Equal(t, "tag1", tags[0].(string))
	assert.Equal(t, "tag2", tags[1].(string))

	// Nested object
	assert.Equal(t, "Test Title", source["object.nestedObject.title"].(string))

	// Multi-dimensional array
	assert.Len(t, source["object.multidimensionalArray"], 2)
	mArr := source["object.multidimensionalArray"].([]interface{})
	assert.ElementsMatch(t, []interface{}{"1", "2", "3"}, mArr[0])
	assert.ElementsMatch(t, []interface{}{"4", "5", "6"}, mArr[1])
}

func Test_NestedObject_Replace(t *testing.T) {

	testSourceSchema := schemer.NewSchema()
	err := schemer.UnmarshalJSON([]byte(testSource), testSourceSchema)
	if err != nil {
		t.Error(err)
	}

	source := testSourceSchema.Normalize(map[string]interface{}{
		"object.nestedObject": map[string]interface{}{
			"title": "Test Title",
		},
	})

	if assert.NotNil(t, source["object.nestedObject"]) {
		o := source["object.nestedObject"].(map[string]interface{})
		assert.Equal(t, "Test Title", o["title"])
	}
}

func Test_NestedObject_ReplaceWithInvalidValue(t *testing.T) {

	testSourceSchema := schemer.NewSchema()
	err := schemer.UnmarshalJSON([]byte(testSource), testSourceSchema)
	if err != nil {
		t.Error(err)
	}

	// Invalid value for map structure
	source := testSourceSchema.Normalize(map[string]interface{}{
		"object.nestedObject":          "Invald Value",
		"object.tags":                  "Invald Value",
		"object.multidimensionalArray": "Invald Value",
	})

	assert.Nil(t, source["object.nestedObject"])
	assert.Nil(t, source["object.tags"])
	assert.Nil(t, source["object.multidimensionalArray"])
}