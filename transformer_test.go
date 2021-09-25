package schemer

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var testSource = `{
	"string": { "type": "string" },
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

var testDest = `{
	"string": { "type": "string" },
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

func TestTransformer(t *testing.T) {

	testSourceSchema := NewSchema()
	err := UnmarsalJSON([]byte(testSource), testSourceSchema)
	if err != nil {
		t.Error(err)
	}

	testDestSchema := NewSchema()
	err = UnmarsalJSON([]byte(testDest), testDestSchema)
	if err != nil {
		t.Error(err)
	}

	// Create transformer
	transformer := NewTransformer(testSourceSchema, testDestSchema)

	// Set transform script
	transformer.SetScript(`
	return {
		string: source.string + 'TEST',
		int: source.int + 1,
		uint: source.uint + 1,
		float: source.float,
		bool: source.bool
	}
`)

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
		t.Error(err)
	}

	returnedValue, err := transformer.Transform(nil, sourceData)
	if err != nil {
		t.Error(err)
	}

	if len(returnedValue) != 1 {
		t.Fail()
	}

	result := returnedValue[0]

	assert.Equal(t, "Brobridge"+"TEST", result["string"].(string))
	assert.Equal(t, int64(-9527)+1, result["int"].(int64))
	assert.Equal(t, uint64(9527)+1, result["uint"].(uint64))
	assert.Equal(t, float64(11.15), result["float"].(float64))
	assert.Equal(t, false, result["bool"].(bool))
}

func TestTransformerEnv(t *testing.T) {

	testSourceSchema := NewSchema()
	err := UnmarsalJSON([]byte(testSource), testSourceSchema)
	if err != nil {
		t.Error(err)
	}

	testDestSchema := NewSchema()
	err = UnmarsalJSON([]byte(testDest), testDestSchema)
	if err != nil {
		t.Error(err)
	}

	// Create transformer
	transformer := NewTransformer(testSourceSchema, testDestSchema)

	// Set transform script
	transformer.SetScript(`
	return {
		string: source.string + env.string
	}
`)

	// Transform
	rawData := `{
	"string": "Brobridge"
}`
	var sourceData map[string]interface{}
	err = json.Unmarshal([]byte(rawData), &sourceData)
	if err != nil {
		t.Error(err)
	}

	returnedValue, err := transformer.Transform(
		map[string]interface{}{
			"string": "test",
		}, sourceData)
	if err != nil {
		t.Error(err)
	}

	if len(returnedValue) != 1 {
		t.Fail()
	}

	result := returnedValue[0]

	assert.Equal(t, "Brobridge"+"test", result["string"].(string))
}

func TestTransformer_MultipleResults(t *testing.T) {

	testSourceSchema := NewSchema()
	err := UnmarsalJSON([]byte(testSource), testSourceSchema)
	if err != nil {
		t.Error(err)
	}

	testDestSchema := NewSchema()
	err = UnmarsalJSON([]byte(testDest), testDestSchema)
	if err != nil {
		t.Error(err)
	}

	// Create transformer
	transformer := NewTransformer(testSourceSchema, testDestSchema)

	// Set transform script
	transformer.SetScript(`
	return [
		{
			string: source.string + 'FIRST'
		},
		{
			string: source.string + 'SECOND'
		}
	]
`)

	// Transform
	rawData := `{
	"string": "Brobridge"
}`
	var sourceData map[string]interface{}
	err = json.Unmarshal([]byte(rawData), &sourceData)
	if err != nil {
		t.Error(err)
	}

	results, err := transformer.Transform(nil, sourceData)
	if err != nil {
		t.Error(err)
	}

	if len(results) != 2 {
		t.Fail()
	}

	assert.Equal(t, "Brobridge"+"FIRST", results[0]["string"].(string))
	assert.Equal(t, "Brobridge"+"SECOND", results[1]["string"].(string))
}

func TestTransformer_Default(t *testing.T) {

	testSourceSchema := NewSchema()
	err := UnmarsalJSON([]byte(testSource), testSourceSchema)
	if err != nil {
		t.Error(err)
	}

	testDestSchema := NewSchema()
	err = UnmarsalJSON([]byte(testDest), testDestSchema)
	if err != nil {
		t.Error(err)
	}

	// Create transformer
	transformer := NewTransformer(testSourceSchema, testDestSchema)

	// Transform
	rawData := `{
	"string": "Brobridge"
}`
	var sourceData map[string]interface{}
	err = json.Unmarshal([]byte(rawData), &sourceData)
	if err != nil {
		t.Error(err)
	}

	results, err := transformer.Transform(nil, sourceData)
	if err != nil {
		t.Error(err)
	}

	if len(results) != 1 {
		t.Fail()
	}

	result := results[0]

	assert.Equal(t, "Brobridge", result["string"].(string))
}

func TestTransformer_NullResult(t *testing.T) {

	testSourceSchema := NewSchema()
	err := UnmarsalJSON([]byte(testSource), testSourceSchema)
	if err != nil {
		t.Error(err)
	}

	testDestSchema := NewSchema()
	err = UnmarsalJSON([]byte(testDest), testDestSchema)
	if err != nil {
		t.Error(err)
	}

	// Create transformer
	transformer := NewTransformer(testSourceSchema, testDestSchema)

	// Set transform script
	transformer.SetScript(`
	return null
`)

	// Transform
	rawData := `{
	"string": "Brobridge"
}`
	var sourceData map[string]interface{}
	err = json.Unmarshal([]byte(rawData), &sourceData)
	if err != nil {
		t.Error(err)
	}

	results, err := transformer.Transform(nil, sourceData)
	if err != nil {
		t.Error(err)
	}

	if len(results) != 0 {
		t.Fail()
	}
}

func TestTransformer_NestedStructure(t *testing.T) {

	sourceSchema := NewSchema()
	err := UnmarsalJSON([]byte(testSource), sourceSchema)
	if err != nil {
		t.Error(err)
	}

	destSchema := NewSchema()
	err = UnmarsalJSON([]byte(testDest), destSchema)
	if err != nil {
		t.Error(err)
	}

	// Create transformer
	transformer := NewTransformer(sourceSchema, destSchema)

	// Set transform script
	transformer.SetScript(`
	return {
		"string": source.object.title,
	}
`)

	// Transform
	rawData := `{
	"string": "Fred Chien",
	"object": {
		"title": "admin",
		"team": "software"
	}
}`
	var sourceData map[string]interface{}
	err = json.Unmarshal([]byte(rawData), &sourceData)
	if err != nil {
		t.Error(err)
	}

	results, err := transformer.Transform(nil, sourceData)
	if err != nil {
		t.Error(err)
	}

	if len(results) != 1 {
		t.Fail()
	}

	result := results[0]

	assert.Equal(t, "admin", result["string"].(string))
}

func TestTransformer_Source_Integer(t *testing.T) {

	sourceSchema := NewSchema()
	err := UnmarsalJSON([]byte(testSource), sourceSchema)
	if err != nil {
		t.Error(err)
	}

	destSchema := NewSchema()
	err = UnmarsalJSON([]byte(testDest), destSchema)
	if err != nil {
		t.Error(err)
	}

	// Create transformer
	transformer := NewTransformer(sourceSchema, destSchema)

	// Set transform script
	transformer.SetScript(`
	return {
		"string": source.int,
		"int": source.int,
		"uint": source.int,
		"float": source.int,
		"bool": source.int
	}
`)

	// Transform
	rawData := `{
	"int": -9527
}`
	var sourceData map[string]interface{}
	err = json.Unmarshal([]byte(rawData), &sourceData)
	if err != nil {
		t.Error(err)
	}

	results, err := transformer.Transform(nil, sourceData)
	if err != nil {
		t.Error(err)
	}

	if len(results) != 1 {
		t.Fail()
	}

	result := results[0]

	assert.Equal(t, "-9527", result["string"].(string))
	assert.Equal(t, int64(-9527), result["int"].(int64))
	assert.Equal(t, uint64(0), result["uint"].(uint64))
	assert.Equal(t, float64(-9527), result["float"].(float64))
	assert.Equal(t, false, result["bool"].(bool))
}

func TestTransformer_Source_UnsignedInteger(t *testing.T) {

	sourceSchema := NewSchema()
	err := UnmarsalJSON([]byte(testSource), sourceSchema)
	if err != nil {
		t.Error(err)
	}

	destSchema := NewSchema()
	err = UnmarsalJSON([]byte(testDest), destSchema)
	if err != nil {
		t.Error(err)
	}

	// Create transformer
	transformer := NewTransformer(sourceSchema, destSchema)

	// Set transform script
	transformer.SetScript(`
	return {
		"string": source.uint,
		"int": source.uint,
		"uint": source.uint,
		"float": source.uint,
		"bool": source.uint
	}
`)

	// Transform
	rawData := `{
	"uint": 9527
}`
	var sourceData map[string]interface{}
	err = json.Unmarshal([]byte(rawData), &sourceData)
	if err != nil {
		t.Error(err)
	}

	results, err := transformer.Transform(nil, sourceData)
	if err != nil {
		t.Error(err)
	}

	if len(results) != 1 {
		t.Fail()
	}

	result := results[0]

	assert.Equal(t, "9527", result["string"].(string))
	assert.Equal(t, int64(9527), result["int"].(int64))
	assert.Equal(t, uint64(9527), result["uint"].(uint64))
	assert.Equal(t, float64(9527), result["float"].(float64))
	assert.Equal(t, true, result["bool"].(bool))
}

func TestTransformer_Source_Float(t *testing.T) {

	sourceSchema := NewSchema()
	err := UnmarsalJSON([]byte(testSource), sourceSchema)
	if err != nil {
		t.Error(err)
	}

	destSchema := NewSchema()
	err = UnmarsalJSON([]byte(testDest), destSchema)
	if err != nil {
		t.Error(err)
	}

	// Create transformer
	transformer := NewTransformer(sourceSchema, destSchema)

	// Set transform script
	transformer.SetScript(`
	return {
		"string": source.float,
		"int": source.float,
		"uint": source.float,
		"float": source.float,
		"bool": source.float
	}
`)

	// Transform
	rawData := `{
	"float": 11.15
}`
	var sourceData map[string]interface{}
	err = json.Unmarshal([]byte(rawData), &sourceData)
	if err != nil {
		t.Error(err)
	}

	results, err := transformer.Transform(nil, sourceData)
	if err != nil {
		t.Error(err)
	}

	if len(results) != 1 {
		t.Fail()
	}

	result := results[0]

	assert.Equal(t, "11.150000", result["string"].(string))
	assert.Equal(t, int64(11), result["int"].(int64))
	assert.Equal(t, uint64(11), result["uint"].(uint64))
	assert.Equal(t, float64(11.15), result["float"].(float64))
	assert.Equal(t, true, result["bool"].(bool))
}

func TestTransformer_Source_String(t *testing.T) {

	sourceSchema := NewSchema()
	err := UnmarsalJSON([]byte(testSource), sourceSchema)
	if err != nil {
		t.Error(err)
	}

	destSchema := NewSchema()
	err = UnmarsalJSON([]byte(testDest), destSchema)
	if err != nil {
		t.Error(err)
	}

	// Create transformer
	transformer := NewTransformer(sourceSchema, destSchema)

	// Set transform script
	transformer.SetScript(`
	return {
		"string": source.string,
		"int": source.string,
		"uint": source.string,
		"float": source.string,
		"bool": source.string,
		"time": "2020-07-19T18:16:08Z",
		"microTime": "2020-07-19T18:16:08Z"
	}
`)

	// Transform
	rawData := `{
	"string": "Brobridge"
}`
	var sourceData map[string]interface{}
	err = json.Unmarshal([]byte(rawData), &sourceData)
	if err != nil {
		t.Error(err)
	}

	results, err := transformer.Transform(nil, sourceData)
	if err != nil {
		t.Error(err)
	}

	if len(results) != 1 {
		t.Fail()
	}

	result := results[0]

	assert.Equal(t, "Brobridge", result["string"].(string))
	assert.Equal(t, int64(0), result["int"].(int64))
	assert.Equal(t, uint64(0), result["uint"].(uint64))
	assert.Equal(t, float64(0), result["float"].(float64))
	assert.Equal(t, false, result["bool"].(bool))
	assert.Equal(t, int64(1595182568), result["time"].(time.Time).Unix())
	assert.Equal(t, int64(1595182568000000000), result["microTime"].(time.Time).UnixNano())
}

func TestTransformer_Source_Bool(t *testing.T) {

	sourceSchema := NewSchema()
	err := UnmarsalJSON([]byte(testSource), sourceSchema)
	if err != nil {
		t.Error(err)
	}

	destSchema := NewSchema()
	err = UnmarsalJSON([]byte(testDest), destSchema)
	if err != nil {
		t.Error(err)
	}

	// Create transformer
	transformer := NewTransformer(sourceSchema, destSchema)

	// Set transform script
	transformer.SetScript(`
	return {
		"string": source.bool,
		"int": source.bool,
		"uint": source.bool,
		"float": source.bool,
		"bool": source.bool
	}
`)

	// Transform
	rawData := `{
	"bool": true
}`
	var sourceData map[string]interface{}
	err = json.Unmarshal([]byte(rawData), &sourceData)
	if err != nil {
		t.Error(err)
	}

	results, err := transformer.Transform(nil, sourceData)
	if err != nil {
		t.Error(err)
	}

	if len(results) != 1 {
		t.Fail()
	}

	result := results[0]

	assert.Equal(t, "true", result["string"].(string))
	assert.Equal(t, int64(1), result["int"].(int64))
	assert.Equal(t, uint64(1), result["uint"].(uint64))
	assert.Equal(t, float64(1), result["float"].(float64))
	assert.Equal(t, true, result["bool"].(bool))
}

func TestTransformer_Source_Time(t *testing.T) {

	sourceSchema := NewSchema()
	err := UnmarsalJSON([]byte(testSource), sourceSchema)
	if err != nil {
		t.Error(err)
	}

	destSchema := NewSchema()
	err = UnmarsalJSON([]byte(testDest), destSchema)
	if err != nil {
		t.Error(err)
	}

	// Create transformer
	transformer := NewTransformer(sourceSchema, destSchema)

	// Set transform script
	transformer.SetScript(`
	return {
		"string": source.time,
		"int": source.time,
		"uint": source.time,
		"float": source.time,
		"bool": source.time,
		"time": source.time,
		"microTime": source.time
	}
`)

	// Transform
	rawData := `{
	"time": 1595182568
}`
	var sourceData map[string]interface{}
	err = json.Unmarshal([]byte(rawData), &sourceData)
	if err != nil {
		t.Error(err)
	}

	results, err := transformer.Transform(nil, sourceData)
	if err != nil {
		t.Error(err)
	}

	if len(results) != 1 {
		t.Fail()
	}

	result := results[0]

	assert.Equal(t, "2020-07-19T18:16:08Z", result["string"].(string))
	assert.Equal(t, int64(1595182568), result["int"].(int64))
	assert.Equal(t, uint64(1595182568), result["uint"].(uint64))
	assert.Equal(t, float64(1595182568), result["float"].(float64))
	assert.Equal(t, true, result["bool"].(bool))
	assert.Equal(t, int64(1595182568), result["time"].(time.Time).Unix())
	assert.Equal(t, int64(1595182568000000000), result["microTime"].(time.Time).UnixNano())
}

func TestTransformer_Source_MicroTime(t *testing.T) {

	sourceSchema := NewSchema()
	err := UnmarsalJSON([]byte(testSource), sourceSchema)
	if err != nil {
		t.Error(err)
	}

	destSchema := NewSchema()
	err = UnmarsalJSON([]byte(testDest), destSchema)
	if err != nil {
		t.Error(err)
	}

	// Create transformer
	transformer := NewTransformer(sourceSchema, destSchema)

	// Set transform script
	transformer.SetScript(`
	return {
		"string": source.microTime,
		"int": source.microTime,
		"uint": source.microTime,
		"float": source.microTime,
		"bool": source.microTime,
		"time": source.microTime,
		"microTime": source.microTime
	}
`)

	// Transform
	rawData := `{
	"microTime": 1595182568000001
}`
	var sourceData map[string]interface{}
	err = json.Unmarshal([]byte(rawData), &sourceData)
	if err != nil {
		t.Error(err)
	}

	results, err := transformer.Transform(nil, sourceData)
	if err != nil {
		t.Error(err)
	}

	if len(results) != 1 {
		t.Fail()
	}

	result := results[0]

	assert.Equal(t, "2020-07-19T18:16:08.000001Z", result["string"].(string))
	assert.Equal(t, int64(1595182568), result["int"].(int64))
	assert.Equal(t, uint64(1595182568), result["uint"].(uint64))
	assert.Equal(t, float64(1595182568), result["float"].(float64))
	assert.Equal(t, true, result["bool"].(bool))
	assert.Equal(t, int64(1595182568), result["time"].(time.Time).Unix())
	assert.Equal(t, int64(1595182568000001000), result["microTime"].(time.Time).UnixNano())
}
