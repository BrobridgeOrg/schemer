package schemer

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var testSource = `{
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

var testDest = `{
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

func TestTransformerScript(t *testing.T) {

	testSourceSchema := NewSchema()
	err := UnmarshalJSON([]byte(testSource), testSourceSchema)
	if err != nil {
		t.Error(err)
	}

	testDestSchema := NewSchema()
	err = UnmarshalJSON([]byte(testDest), testDestSchema)
	if err != nil {
		t.Error(err)
	}

	// Create transformer
	transformer := NewTransformer(testSourceSchema, testDestSchema)

	// Set transform script
	err = transformer.SetScript(`
		Invalid script
`)
	//	t.Log(err)
	assert.NotNil(t, err)
}

func TestTransformer(t *testing.T) {

	testSourceSchema := NewSchema()
	err := UnmarshalJSON([]byte(testSource), testSourceSchema)
	if err != nil {
		t.Error(err)
	}

	testDestSchema := NewSchema()
	err = UnmarshalJSON([]byte(testDest), testDestSchema)
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
	if !assert.Nil(t, err) {
		return
	}

	returnedValue, err := transformer.Transform(nil, sourceData)
	if !assert.Nil(t, err) {
		return
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

func TestTransformerWithoutSourceSchema(t *testing.T) {

	testDestSchema := NewSchema()
	err := UnmarshalJSON([]byte(testDest), testDestSchema)
	if err != nil {
		t.Error(err)
	}

	// Create transformer
	transformer := NewTransformer(nil, testDestSchema)

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
	if !assert.Nil(t, err) {
		return
	}

	returnedValue, err := transformer.Transform(nil, sourceData)
	if !assert.Nil(t, err) {
		return
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

func TestTransformerWithoutSchema(t *testing.T) {

	// Create transformer
	transformer := NewTransformer(nil, nil)

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
	err := json.Unmarshal([]byte(rawData), &sourceData)
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
	assert.Equal(t, int64(9527)+1, result["uint"].(int64))
	assert.Equal(t, float64(11.15), result["float"].(float64))
	assert.Equal(t, false, result["bool"].(bool))
}

func TestTransformerEnv(t *testing.T) {

	testSourceSchema := NewSchema()
	err := UnmarshalJSON([]byte(testSource), testSourceSchema)
	if err != nil {
		t.Error(err)
	}

	testDestSchema := NewSchema()
	err = UnmarshalJSON([]byte(testDest), testDestSchema)
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
	err := UnmarshalJSON([]byte(testSource), testSourceSchema)
	if err != nil {
		t.Error(err)
	}

	testDestSchema := NewSchema()
	err = UnmarshalJSON([]byte(testDest), testDestSchema)
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
	err := UnmarshalJSON([]byte(testSource), testSourceSchema)
	if err != nil {
		t.Error(err)
	}

	testDestSchema := NewSchema()
	err = UnmarshalJSON([]byte(testDest), testDestSchema)
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
	if !assert.Nil(t, err) {
		return
	}

	if len(results) != 1 {
		t.Fail()
	}

	result := results[0]

	assert.Equal(t, "Brobridge", result["string"].(string))
}

func TestTransformer_NullResult(t *testing.T) {

	testSourceSchema := NewSchema()
	err := UnmarshalJSON([]byte(testSource), testSourceSchema)
	if err != nil {
		t.Error(err)
	}

	testDestSchema := NewSchema()
	err = UnmarshalJSON([]byte(testDest), testDestSchema)
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

	assert.Equal(t, 0, len(results))
}

func TestTransformer_NestedStructure(t *testing.T) {

	sourceSchema := NewSchema()
	err := UnmarshalJSON([]byte(testSource), sourceSchema)
	if err != nil {
		t.Error(err)
	}

	destSchema := NewSchema()
	err = UnmarshalJSON([]byte(testDest), destSchema)
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

func TestTransformer_Source_Binary(t *testing.T) {

	sourceSchema := NewSchema()
	err := UnmarshalJSON([]byte(testSource), sourceSchema)
	if err != nil {
		t.Error(err)
	}

	destSchema := NewSchema()
	err = UnmarshalJSON([]byte(testDest), destSchema)
	if err != nil {
		t.Error(err)
	}

	// Create transformer
	transformer := NewTransformer(sourceSchema, destSchema)

	// Set transform script
	transformer.SetScript(`
	return {
		"string": source.binary,
		"binary": source.binary,
		"int": source.binary,
		"uint": source.binary,
		"float": source.binary,
		"bool": source.binary
	}
`)

	// Transform
	rawData := `{
	"binary": [ 0, 1, 2, 3, 4, 5, 6, 7, 8, 9 ]
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

	assert.Equal(t, "[0 1 2 3 4 5 6 7 8 9]", result["string"].(string))
	assert.Equal(t, []uint8{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, result["binary"])
	assert.Equal(t, int64(0), result["int"].(int64))
	assert.Equal(t, uint64(0), result["uint"].(uint64))
	assert.Equal(t, float64(0), result["float"].(float64))
	assert.Equal(t, false, result["bool"].(bool))
}

func TestTransformer_Source_Integer(t *testing.T) {

	sourceSchema := NewSchema()
	err := UnmarshalJSON([]byte(testSource), sourceSchema)
	if err != nil {
		t.Error(err)
	}

	destSchema := NewSchema()
	err = UnmarshalJSON([]byte(testDest), destSchema)
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
	err := UnmarshalJSON([]byte(testSource), sourceSchema)
	if err != nil {
		t.Error(err)
	}

	destSchema := NewSchema()
	err = UnmarshalJSON([]byte(testDest), destSchema)
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
	err := UnmarshalJSON([]byte(testSource), sourceSchema)
	if err != nil {
		t.Error(err)
	}

	destSchema := NewSchema()
	err = UnmarshalJSON([]byte(testDest), destSchema)
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

	assert.Equal(t, "11.15", result["string"].(string))
	assert.Equal(t, int64(11), result["int"].(int64))
	assert.Equal(t, uint64(11), result["uint"].(uint64))
	assert.Equal(t, float64(11.15), result["float"].(float64))
	assert.Equal(t, true, result["bool"].(bool))
}

func TestTransformer_Source_String(t *testing.T) {

	sourceSchema := NewSchema()
	err := UnmarshalJSON([]byte(testSource), sourceSchema)
	if err != nil {
		t.Error(err)
	}

	destSchema := NewSchema()
	err = UnmarshalJSON([]byte(testDest), destSchema)
	if err != nil {
		t.Error(err)
	}

	// Create transformer
	transformer := NewTransformer(sourceSchema, destSchema)

	// Set transform script
	transformer.SetScript(`
	return {
		"string": source.string,
		"binary": source.string,
		"int": source.string,
		"uint": source.string,
		"float": source.string,
		"bool": source.string,
		"time": "2020-07-19T18:16:08Z",
		"microTime": "2020-07-19 18:16:08.1234567"
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
	assert.Equal(t, []byte("Brobridge"), result["binary"])
	assert.Equal(t, int64(0), result["int"].(int64))
	assert.Equal(t, uint64(0), result["uint"].(uint64))
	assert.Equal(t, float64(0), result["float"].(float64))
	assert.Equal(t, false, result["bool"].(bool))
	assert.Equal(t, int64(1595182568), result["time"].(time.Time).Unix())
	assert.Equal(t, int64(1595182568123456700), result["microTime"].(time.Time).UnixNano())
}

func TestTransformer_Source_Bool_With_True(t *testing.T) {

	sourceSchema := NewSchema()
	err := UnmarshalJSON([]byte(testSource), sourceSchema)
	if err != nil {
		t.Error(err)
	}

	destSchema := NewSchema()
	err = UnmarshalJSON([]byte(testDest), destSchema)
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

func TestTransformer_Source_Bool_With_False(t *testing.T) {

	sourceSchema := NewSchema()
	err := UnmarshalJSON([]byte(testSource), sourceSchema)
	if err != nil {
		t.Error(err)
	}

	destSchema := NewSchema()
	err = UnmarshalJSON([]byte(testDest), destSchema)
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
	"bool": false
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

	assert.Equal(t, "false", result["string"].(string))
	assert.Equal(t, int64(0), result["int"].(int64))
	assert.Equal(t, uint64(0), result["uint"].(uint64))
	assert.Equal(t, float64(0), result["float"].(float64))
	assert.Equal(t, false, result["bool"].(bool))
}

func TestTransformer_Source_Time(t *testing.T) {

	sourceSchema := NewSchema()
	err := UnmarshalJSON([]byte(testSource), sourceSchema)
	if err != nil {
		t.Error(err)
	}

	destSchema := NewSchema()
	err = UnmarshalJSON([]byte(testDest), destSchema)
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

func TestTransformer_Source_Time_Dest_Empty(t *testing.T) {

	sourceSchema := NewSchema()
	err := UnmarshalJSON([]byte(testSource), sourceSchema)
	if err != nil {
		t.Error(err)
	}

	// Create transformer
	transformer := NewTransformer(sourceSchema, nil)

	// Set transform script
	transformer.SetScript(`
	return source;
`)

	// Transform
	rawData := `{
	"time": 1595182568
}`
	var sourceData map[string]interface{}
	err = json.Unmarshal([]byte(rawData), &sourceData)
	if !assert.Nil(t, err) {
		return
	}

	results, err := transformer.Transform(nil, sourceData)
	if !assert.Nil(t, err) {
		return
	}

	if len(results) != 1 {
		t.Fail()
	}

	result := results[0]

	assert.Equal(t, int64(1595182568), result["time"].(time.Time).Unix())
}

func TestTransformer_Source_MicroTime(t *testing.T) {

	sourceSchema := NewSchema()
	err := UnmarshalJSON([]byte(testSource), sourceSchema)
	if !assert.Nil(t, err) {
		return
	}

	destSchema := NewSchema()
	err = UnmarshalJSON([]byte(testDest), destSchema)
	if !assert.Nil(t, err) {
		return
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

func TestTransformer_Source_Null(t *testing.T) {

	sourceSchema := NewSchema()
	err := UnmarshalJSON([]byte(testSource), sourceSchema)
	if err != nil {
		t.Error(err)
	}

	destSchema := NewSchema()
	err = UnmarshalJSON([]byte(testDest), destSchema)
	if err != nil {
		t.Error(err)
	}

	// Create transformer
	transformer := NewTransformer(sourceSchema, destSchema)

	// Set transform script
	transformer.SetScript(`
	return {
		"string": source.string,
		"binary": source.string,
		"int": source.string,
		"uint": source.string,
		"float": source.string,
		"bool": source.string,
		"time": source.string,
		"microTime": source.string,
	}
`)

	// Transform
	rawData := `{
	"string": null
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

	assert.Equal(t, nil, result["string"])
	assert.Equal(t, nil, result["binary"])
	assert.Equal(t, nil, result["int"])
	assert.Equal(t, nil, result["uint"])
	assert.Equal(t, nil, result["float"])
	assert.Equal(t, nil, result["bool"])
	assert.Equal(t, nil, result["time"])
	assert.Equal(t, nil, result["microTime"])
}

func TestTransformer_Source_TimeString(t *testing.T) {

	sourceSchema := NewSchema()
	err := UnmarshalJSON([]byte(testSource), sourceSchema)
	if err != nil {
		t.Error(err)
	}

	destSchema := NewSchema()
	err = UnmarshalJSON([]byte(testDest), destSchema)
	if err != nil {
		t.Error(err)
	}

	// Create transformer
	transformer := NewTransformer(sourceSchema, destSchema)

	// Set transform script
	transformer.SetScript(`
	return {
		"string": source.string,
		"binary": source.string,
		"int": source.string,
		"uint": source.string,
		"float": source.string,
		"bool": source.string,
		"time": source.time,
		"microTime": "2020-07-19 18:16:08.1234567"
	}
`)

	// Transform
	rawData := `{
	"string": "Brobridge",
	"time": "2020-07-19T18:16:08.96Z"
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
	assert.Equal(t, []byte("Brobridge"), result["binary"])
	assert.Equal(t, int64(0), result["int"].(int64))
	assert.Equal(t, uint64(0), result["uint"].(uint64))
	assert.Equal(t, float64(0), result["float"].(float64))
	assert.Equal(t, false, result["bool"].(bool))
	assert.Equal(t, int64(1595182568), result["time"].(time.Time).Unix())
	assert.Equal(t, int64(1595182568123456700), result["microTime"].(time.Time).UnixNano())
}

func TestTransformer_Source_TimeEmptyString(t *testing.T) {

	sourceSchema := NewSchema()
	err := UnmarshalJSON([]byte(testSource), sourceSchema)
	if err != nil {
		t.Error(err)
	}

	destSchema := NewSchema()
	err = UnmarshalJSON([]byte(testDest), destSchema)
	if err != nil {
		t.Error(err)
	}

	// Create transformer
	transformer := NewTransformer(sourceSchema, destSchema)

	// Set transform script
	transformer.SetScript(`
	return {
		"time": source.time,
	}
`)

	// Transform
	rawData := `{
	"time": ""
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

	assert.Equal(t, nil, result["time"])
}
