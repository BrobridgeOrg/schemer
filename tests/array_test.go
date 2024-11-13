package tests

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/BrobridgeOrg/schemer"
	"github.com/stretchr/testify/assert"
)

var arraySchema1 = `{
    "array_string":{
        "type":"array",
        "subtype":"string"
    },
    "array_binary":{
        "type":"array",
        "subtype":"binary"
    },
    "array_int":{
        "type":"array",
        "subtype":"int"
    },
    "array_uint":{
        "type":"array",
        "subtype":"uint"
    },
    "array_float":{
        "type":"array",
        "subtype":"float"
    },
    "array_bool":{
        "type":"array",
        "subtype":"bool"
    },
    "array_time":{
        "type":"array",
        "subtype":"time"
    },
    "array_any":{
        "type":"array",
        "subtype":"any"
    }
}`

func TestEmptyArray(t *testing.T) {

	transformer, testSourceSchema := SetupTransformer(t, arraySchema1)
	var rawData map[string]interface{}
	err := json.Unmarshal([]byte(`{"array_string":[], "array_int":[], "array_uint":[], "array_binary":[], "array_float":[], "array_bool":[], "array_time":[], "array_any":[]}`), &rawData)
	if err != nil {
		t.Fatal(err)
	}

	source := testSourceSchema.Normalize(rawData)
	result := AssertTransform(t, transformer, source)

	assert.Equal(t, []interface{}{}, result["array_string"])
	assert.Equal(t, []interface{}{}, result["array_int"])
	assert.Equal(t, []interface{}{}, result["array_uint"])
	assert.Equal(t, []interface{}{}, result["array_binary"])
	assert.Equal(t, []interface{}{}, result["array_float"])
	assert.Equal(t, []interface{}{}, result["array_bool"])
	assert.Equal(t, []interface{}{}, result["array_time"])
	assert.Equal(t, []interface{}{}, result["array_any"])
}

func TestStringArray(t *testing.T) {

	transformer, testSourceSchema := SetupTransformer(t, arraySchema1)
	var rawData map[string]interface{}
	err := json.Unmarshal([]byte(`{"array_string":["a", "b", "c"], "array_binary":["a", "b", "c"], "array_any":["a", "b", "c"]}`), &rawData)
	if err != nil {
		t.Fatal(err)
	}

	source := testSourceSchema.Normalize(rawData)
	result := AssertTransform(t, transformer, source)

	assert.ElementsMatch(t, []string{"a", "b", "c"}, result["array_string"])
	assert.ElementsMatch(t, [][]byte{{0x61}, {0x62}, {0x63}}, result["array_binary"])
	assert.ElementsMatch(t, []string{"a", "b", "c"}, result["array_any"])
}

func TestIntArray(t *testing.T) {

	transformer, testSourceSchema := SetupTransformer(t, arraySchema1)
	var rawData map[string]interface{}
	err := json.Unmarshal([]byte(`{"array_int":[1, 2, 3], "array_uint":[1, 2, 3], "array_float":[1, 2, 3], "array_any":[1, 2, 3]}`), &rawData)
	if err != nil {
		t.Fatal(err)
	}

	source := testSourceSchema.Normalize(rawData)
	result := AssertTransform(t, transformer, source)

	assert.ElementsMatch(t, []int64{1, 2, 3}, result["array_int"])
	assert.ElementsMatch(t, []uint64{1, 2, 3}, result["array_uint"])
	assert.ElementsMatch(t, []float64{1, 2, 3}, result["array_float"])
	assert.ElementsMatch(t, []float64{1, 2, 3}, result["array_any"])
}

func TestIntArray2(t *testing.T) {

	transformer, testSourceSchema := SetupTransformer(t, arraySchema1)
	var rawData map[string]interface{}
	err := json.Unmarshal([]byte(`{"array_int":[1, "2", 3], "array_uint":[1, "2", 3], "array_float":[1, "2", 3], "array_any":[1, "2", 3]}`), &rawData)
	if err != nil {
		t.Fatal(err)
	}

	source := testSourceSchema.Normalize(rawData)
	result := AssertTransform(t, transformer, source)

	assert.ElementsMatch(t, []int64{1, 2, 3}, result["array_int"])
	assert.ElementsMatch(t, []uint64{1, 2, 3}, result["array_uint"])
	assert.ElementsMatch(t, []float64{1, 2, 3}, result["array_float"])
	assert.ElementsMatch(t, []interface{}{float64(1), "2", float64(3)}, result["array_any"])
}

func TestBinaryArray(t *testing.T) {

	transformer, testSourceSchema := SetupTransformer(t, arraySchema1)
	var rawData map[string]interface{}
	err := json.Unmarshal([]byte(`{"array_binary":["00", "01", "10", "11"]}`), &rawData)
	if err != nil {
		t.Fatal(err)
	}

	source := testSourceSchema.Normalize(rawData)
	result := AssertTransform(t, transformer, source)

	assert.ElementsMatch(t, [][]byte{{0x30, 0x30}, {0x30, 0x31}, {0x31, 0x30}, {0x31, 0x31}}, result["array_binary"])
}

func TestFloatArray(t *testing.T) {

	transformer, testSourceSchema := SetupTransformer(t, arraySchema1)
	var rawData map[string]interface{}
	err := json.Unmarshal([]byte(`{"array_float":[1.1, 2.2, 3.3], "array_any":[1.1, 2.2, 3.3]}`), &rawData)
	if err != nil {
		t.Fatal(err)
	}

	source := testSourceSchema.Normalize(rawData)
	result := AssertTransform(t, transformer, source)

	assert.ElementsMatch(t, []float64{1.1, 2.2, 3.3}, result["array_float"])
	assert.ElementsMatch(t, []float64{1.1, 2.2, 3.3}, result["array_any"])
}

func TestBoolArray(t *testing.T) {

	transformer, testSourceSchema := SetupTransformer(t, arraySchema1)
	var rawData map[string]interface{}
	err := json.Unmarshal([]byte(`{"array_bool":[true, false], "array_any":[true, false]}`), &rawData)
	if err != nil {
		t.Fatal(err)
	}

	source := testSourceSchema.Normalize(rawData)
	result := AssertTransform(t, transformer, source)

	assert.ElementsMatch(t, []bool{true, false}, result["array_bool"])
	assert.ElementsMatch(t, []bool{true, false}, result["array_any"])
}

func TestTimeArray(t *testing.T) {

	transformer, testSourceSchema := SetupTransformer(t, arraySchema1)
	var rawData map[string]interface{}
	err := json.Unmarshal([]byte(`{"array_time":["2024-08-06T15:02:00Z", "2024-08-06T15:03:00Z", "2024-08-06T15:04:00Z"], "array_any":["2024-08-06T15:02:00Z", "2024-08-06T15:02:00Z", "2024-08-06T15:02:00Z"]}`), &rawData)
	if err != nil {
		t.Fatal(err)
	}

	source := testSourceSchema.Normalize(rawData)
	result := AssertTransform(t, transformer, source)

	resultTimeArray := []time.Time{}
	for _, t := range result["array_time"].([]interface{}) {
		resultTimeArray = append(resultTimeArray, t.(time.Time).UTC())
	}

	assert.ElementsMatch(t, []time.Time{time.Date(2024, 8, 6, 15, 2, 0, 0, time.UTC), time.Date(2024, 8, 6, 15, 3, 0, 0, time.UTC), time.Date(2024, 8, 6, 15, 4, 0, 0, time.UTC)}, resultTimeArray)
	assert.ElementsMatch(t, []string{"2024-08-06T15:02:00Z", "2024-08-06T15:02:00Z", "2024-08-06T15:02:00Z"}, result["array_any"])
}

func TestLongStringArray(t *testing.T) {

	transformer, testSourceSchema := SetupTransformer(t, arraySchema1)
	longString := strings.Repeat("a", 32768)
	testJSON := fmt.Sprintf(`{"array_string":["%s", "%s", "%s"], "array_any":["%s", "%s", "%s"]}`, longString, longString, longString, longString, longString, longString)
	var rawData map[string]interface{}
	err := json.Unmarshal([]byte(testJSON), &rawData)
	if err != nil {
		t.Fatal(err)
	}

	source := testSourceSchema.Normalize(rawData)
	result := AssertTransform(t, transformer, source)
	
	assert.ElementsMatch(t, []string{longString, longString, longString}, result["array_string"])
	assert.ElementsMatch(t, []string{longString, longString, longString}, result["array_any"])
}

func TestSingleElementArray(t *testing.T) {

	transformer, testSourceSchema := SetupTransformer(t, arraySchema1)
	var rawData map[string]interface{}
	err := json.Unmarshal([]byte(`{"array_int":[1], "array_string":["a"], "array_any":[1]}`), &rawData)
	if err != nil {
		t.Fatal(err)
	}

	source := testSourceSchema.Normalize(rawData)
	result := AssertTransform(t, transformer, source)

	assert.ElementsMatch(t, []int64{1}, result["array_int"])
	assert.ElementsMatch(t, []string{"a"}, result["array_string"])
	assert.ElementsMatch(t, []float64{1}, result["array_any"])
}

func TestMassiveElementsArray(t *testing.T) {

	transformer, testSourceSchema := SetupTransformer(t, arraySchema1)
	massiveElement := ""
	for i := 1; i <= 32768; i++ {
		if i > 1 {
			massiveElement += ", "
		}
		massiveElement += fmt.Sprintf("%d", i)
	}
	testJSON := fmt.Sprintf(`{"array_int":[%s], "array_any":[%s]}`, massiveElement, massiveElement)

	var rawData map[string]interface{}
	err := json.Unmarshal([]byte(testJSON), &rawData)
	if err != nil {
		t.Fatal(err)
	}

	source := testSourceSchema.Normalize(rawData)
	result := AssertTransform(t, transformer, source)
	expectedInt := make([]int64, 32768)
	expectedAny := make([]float64, 32768)
	for i := 1; i <= 32768; i++ {
		expectedInt[i-1] = int64(i)
		expectedAny[i-1] = float64(i)
	}

	assert.ElementsMatch(t, expectedInt, result["array_int"])
	assert.ElementsMatch(t, expectedAny, result["array_any"])
}

func TestSameElementArray(t *testing.T) {

	transformer, testSourceSchema := SetupTransformer(t, arraySchema1)
	var rawData map[string]interface{}
	err := json.Unmarshal([]byte(`{"array_string":["a", "a", "a"], "array_int":[1, 1, 1], "array_any":[1, 1, 1]}`), &rawData)
	if err != nil {
		t.Fatal(err)
	}

	source := testSourceSchema.Normalize(rawData)
	result := AssertTransform(t, transformer, source)

	assert.ElementsMatch(t, []string{"a", "a", "a"}, result["array_string"])
	assert.ElementsMatch(t, []int64{1, 1, 1}, result["array_int"])
	assert.ElementsMatch(t, []float64{1, 1, 1}, result["array_any"])
}

func TestStringArrayWithInvalidSubtypes(t *testing.T) {

	transformer, testSourceSchema := SetupTransformer(t, arraySchema1)
	var rawData map[string]interface{}
	err := json.Unmarshal([]byte(`{"array_int":["a", "b", "c"], "array_uint":["a", "b", "c"], "array_float":["a", "b", "c"], "array_bool":["a", "b", "c"], "array_time":["a", "b", "c"]}`), &rawData)
	if err != nil {
		t.Fatal(err)
	}

	source := testSourceSchema.Normalize(rawData)
	result := AssertTransform(t, transformer, source)

	assert.ElementsMatch(t, nil, result["array_int"])
	assert.ElementsMatch(t, nil, result["array_uint"])
	assert.ElementsMatch(t, nil, result["array_float"])
	assert.ElementsMatch(t, nil, result["array_bool"])
	var tt time.Time
	assert.ElementsMatch(t, []time.Time{tt, tt, tt}, result["array_time"])
}

func TestIntArrayWithInvalidSubypes(t *testing.T) {

	transformer, testSourceSchema := SetupTransformer(t, arraySchema1)
	var rawData map[string]interface{}
	err := json.Unmarshal([]byte(`{"array_string":[1, 2, 3], "array_binary":[1, 2, 3], "array_bool":[1, 2, 3], "array_time":[1, 2, 3]}`), &rawData)
	if err != nil {
		t.Fatal(err)
	}

	source := testSourceSchema.Normalize(rawData)
	result := AssertTransform(t, transformer, source)

	assert.ElementsMatch(t, []string{"1", "2", "3"}, result["array_string"])
	assert.ElementsMatch(t, nil, result["array_binary"])
	assert.ElementsMatch(t, []bool{true, true, true}, result["array_bool"])
	assert.ElementsMatch(t, []time.Time{time.Unix(1, 0), time.Unix(2, 0), time.Unix(3, 0)}, result["array_time"])
}

func TestBinaryArrayWithInvalidSubtypes(t *testing.T) {

	transformer, testSourceSchema := SetupTransformer(t, arraySchema1)
	var rawData map[string]interface{}
	err := json.Unmarshal([]byte(`{"array_int":["00", "01", "10", "11"], "array_uint":["00", "01", "10", "11"], "array_float":["00", "01", "10", "11"], "array_bool":["00", "01", "10", "11"], "array_time":["00", "01", "10", "11"]}`), &rawData)
	if err != nil {
		t.Fatal(err)
	}

	source := testSourceSchema.Normalize(rawData)
	result := AssertTransform(t, transformer, source)

	assert.ElementsMatch(t, []int64{0, 1, 10, 11}, result["array_int"])
	assert.ElementsMatch(t, []uint64{0, 1, 10, 11}, result["array_uint"])
	assert.ElementsMatch(t, []float64{0, 1, 10, 11}, result["array_float"])
	assert.ElementsMatch(t, nil, result["array_bool"])
	var tt time.Time
	assert.ElementsMatch(t, []time.Time{tt, tt, tt, tt}, result["array_time"])
}

func TestFloatArrayWithInvalidSubtypes(t *testing.T) {

	transformer, testSourceSchema := SetupTransformer(t, arraySchema1)
	var rawData map[string]interface{}
	err := json.Unmarshal([]byte(`{"array_string":[1.1, 2.2, 3.3], "array_int":[1.1, 2.2, 3.3], "array_uint":[1.1, 2.2, 3.3], "array_binary":[1.1, 2.2, 3.3], "array_bool":[1.1, 2.2, 3.3], "array_time":[1.1, 2.2, 3.3]}`), &rawData)
	if err != nil {
		t.Fatal(err)
	}

	source := testSourceSchema.Normalize(rawData)
	result := AssertTransform(t, transformer, source)

	assert.ElementsMatch(t, []string{"1.1", "2.2", "3.3"}, result["array_string"])
	assert.ElementsMatch(t, []int64{1, 2, 3}, result["array_int"])
	assert.ElementsMatch(t, []uint64{1, 2, 3}, result["array_uint"])
	assert.ElementsMatch(t, nil, result["array_binary"])
	assert.ElementsMatch(t, []bool{true, true, true}, result["array_bool"])
	assert.ElementsMatch(t, []time.Time{time.Unix(1, 0), time.Unix(2, 0), time.Unix(3, 0)}, result["array_time"])
}

func TestMixTypesArray(t *testing.T) {

	transformer, testSourceSchema := SetupTransformer(t, arraySchema1)
	var rawData map[string]interface{}
	err := json.Unmarshal([]byte(`{"array_string":[1, "a", 5.6], "array_int":[1, "a", 5.6], "array_uint":[1, "a", 5.6], "array_binary":[1, "a", 5.6], "array_float":[1, "a", 5.6], "array_bool":[1, "a", 5.6], "array_time":[1, "a", 5.6], "array_any":[1, "a", 5.6, true]}`), &rawData)
	if err != nil {
		t.Fatal(err)
	}

	source := testSourceSchema.Normalize(rawData)
	result := AssertTransform(t, transformer, source)

	assert.ElementsMatch(t, []string{"1", "a", "5.6"}, result["array_string"])
	assert.ElementsMatch(t, nil, result["array_int"])
	assert.ElementsMatch(t, nil, result["array_uint"])
	assert.ElementsMatch(t, nil, result["array_binary"])
	assert.ElementsMatch(t, nil, result["array_float"])
	assert.ElementsMatch(t, nil, result["array_bool"])
	var tt time.Time
	assert.ElementsMatch(t, []time.Time{time.Unix(1, 0), tt, time.Unix(5, 0)}, result["array_time"])
	assert.ElementsMatch(t, []interface{}{float64(1), "a", 5.6, true}, result["array_any"])
}

func TestNonArray_Null(t *testing.T) {

	transformer, testSourceSchema := SetupTransformer(t, arraySchema1)
	var rawData map[string]interface{}
	err := json.Unmarshal([]byte(`{"array_string":"", "array_int":"", "array_uint":"", "array_binary":"", "array_float":"", "array_bool":"", "array_time":"", "array_any":""}`), &rawData)
	if err != nil {
		t.Fatal(err)
	}

	source := testSourceSchema.Normalize(rawData)
	result := AssertTransform(t, transformer, source)

	assert.ElementsMatch(t, nil, result["array_string"])
	assert.ElementsMatch(t, nil, result["array_int"])
	assert.ElementsMatch(t, nil, result["array_uint"])
	assert.ElementsMatch(t, nil, result["array_binary"])
	assert.ElementsMatch(t, nil, result["array_float"])
	assert.ElementsMatch(t, nil, result["array_bool"])
	assert.ElementsMatch(t, nil, result["array_time"])
	assert.ElementsMatch(t, nil, result["array_any"])
}

func TestNonArray_Space(t *testing.T) {

	transformer, testSourceSchema := SetupTransformer(t, arraySchema1)
	var rawData map[string]interface{}
	err := json.Unmarshal([]byte(`{"array_string":" ", "array_int":" ", "array_uint":" ", "array_binary":" ", "array_float":" ", "array_bool":" ", "array_time":" ", "array_any":" "}`), &rawData)
	if err != nil {
		t.Fatal(err)
	}

	source := testSourceSchema.Normalize(rawData)
	result := AssertTransform(t, transformer, source)

	assert.ElementsMatch(t, nil, result["array_string"])
	assert.ElementsMatch(t, nil, result["array_int"])
	assert.ElementsMatch(t, nil, result["array_uint"])
	assert.ElementsMatch(t, nil, result["array_binary"])
	assert.ElementsMatch(t, nil, result["array_float"])
	assert.ElementsMatch(t, nil, result["array_bool"])
	assert.ElementsMatch(t, nil, result["array_time"])
	assert.ElementsMatch(t, nil, result["array_any"])
}

func TestNonArray_String(t *testing.T) {

	transformer, testSourceSchema := SetupTransformer(t, arraySchema1)

	var rawData map[string]interface{}
	err := json.Unmarshal([]byte(`{"array_string":"abc", "array_int":"abc", "array_uint":"abc", "array_binary":"abc", "array_float":"abc", "array_bool":"abc", "array_time":"abc", "array_any":"abc"}`), &rawData)
	if err != nil {
		t.Fatal(err)
	}

	source := testSourceSchema.Normalize(rawData)
	result := AssertTransform(t, transformer, source)

	assert.ElementsMatch(t, nil, result["array_string"])
	assert.ElementsMatch(t, nil, result["array_int"])
	assert.ElementsMatch(t, nil, result["array_uint"])
	assert.ElementsMatch(t, nil, result["array_binary"])
	assert.ElementsMatch(t, nil, result["array_float"])
	assert.ElementsMatch(t, nil, result["array_bool"])
	assert.ElementsMatch(t, nil, result["array_time"])
	assert.ElementsMatch(t, nil, result["array_any"])
}

func TestNonArray_Int(t *testing.T) {

	transformer, testSourceSchema := SetupTransformer(t, arraySchema1)

	var rawData map[string]interface{}
	err := json.Unmarshal([]byte(`{"array_string":5, "array_int":5, "array_uint":5, "array_binary":5, "array_float":5, "array_bool":5, "array_time":5, "array_any":5}`), &rawData)
	if err != nil {
		t.Fatal(err)
	}

	source := testSourceSchema.Normalize(rawData)
	result := AssertTransform(t, transformer, source)

	assert.ElementsMatch(t, nil, result["array_string"])
	assert.ElementsMatch(t, nil, result["array_int"])
	assert.ElementsMatch(t, nil, result["array_uint"])
	assert.ElementsMatch(t, nil, result["array_binary"])
	assert.ElementsMatch(t, nil, result["array_float"])
	assert.ElementsMatch(t, nil, result["array_bool"])
	assert.ElementsMatch(t, nil, result["array_time"])
	assert.ElementsMatch(t, nil, result["array_any"])
}

func SetupTransformer(t *testing.T, schema string) (*schemer.Transformer, *schemer.Schema) {

	testSourceSchema := schemer.NewSchema()
	err := schemer.UnmarshalJSON([]byte(schema), testSourceSchema)
	if err != nil {
		t.Error(err)
	}

	// Using the same schema for destination
	testDestSchema := schemer.NewSchema()
	err = schemer.UnmarshalJSON([]byte(schema), testDestSchema)
	if err != nil {
		t.Error(err)
	}

	// Create transformer
	transformer := schemer.NewTransformer(testSourceSchema, testDestSchema)
	err = transformer.SetScript(`return source`)
	if err != nil {
		t.Error(err)
	}

	return transformer, testSourceSchema
}

func AssertTransform(t *testing.T, transformer *schemer.Transformer, source map[string]interface{}) map[string]interface{} {

	returnedValue, err := transformer.Transform(nil, source)
	if !assert.Nil(t, err) {
		t.Fatal(err)
	}

	if !assert.Len(t, returnedValue, 1) {
		t.Fatal(err)
	}

	result := returnedValue[0]
	return result
}
