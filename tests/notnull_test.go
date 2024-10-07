package tests

import (
	"testing"

	"github.com/BrobridgeOrg/schemer"
	"github.com/stretchr/testify/assert"
)

var notNullSchema = `{
	"bool": {
		"type": "bool",
		"notNull": true
	}
}`

func Test_NotNull(t *testing.T) {

	var schema = `{
		"bool_true": {
			"type": "bool",
			"notNull": true
		},
		"bool_false": {
			"type": "bool",
			"notNull": true
		},
		"bool_null": {
			"type": "bool",
			"notNull": true
		}
	}`

	s := schemer.NewSchema()
	err := schemer.UnmarshalJSON([]byte(schema), s)
	if err != nil {
		t.Error(err)
	}

	source := s.Normalize(map[string]interface{}{
		"bool_true":  true,
		"bool_false": false,
		"bool_null":  nil,
	})

	assert.Equal(t, true, source["bool_true"])
	assert.Equal(t, false, source["bool_false"])
	assert.Equal(t, false, source["bool_null"])

	// Create transformer
	transformer := schemer.NewTransformer(s, s)

	// Set transform script
	transformer.SetScript(`return source`)

	// Transforming
	returnedValue, err := transformer.Transform(nil, source)
	if !assert.Nil(t, err) {
		return
	}

	if !assert.Len(t, returnedValue, 1) {
		return
	}

	result := returnedValue[0]

	assert.Equal(t, true, result["bool_true"])
	assert.Equal(t, false, result["bool_false"])
	assert.Equal(t, false, result["bool_null"])
}
