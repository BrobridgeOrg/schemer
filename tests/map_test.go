package tests

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/BrobridgeOrg/schemer"
	"github.com/stretchr/testify/assert"
)

var mapSchema1 = `
{
   "id":{
      "type":"uint"
   },
   "map_col":{
      "type":"map"
   }
}
`

var mapSchema2 = `
{
   "id":{
      "type":"uint"
   },
   "map_col":{
      "type":"map",
      "fields":{
         "time_col":{
            "type":"time",
            "precision":"milisecond"
         }
      }
   }
}
`

var mapSchema3 = `
{
   "id":{
      "type":"uint"
   },
   "map_col":{
      "type":"map",
      "fields":{
         "string_col":{
            "type":"string"
         },
         "binary_col":{
            "type":"binary"
         },
         "int_col":{
            "type":"int"
         },
         "uint_col":{
            "type":"uint"
         },
         "float_col":{
            "type":"float"
         },
         "bool_col":{
            "type":"bool"
         },
         "any_col":{
            "type":"any"
         }
      }
   }
}
`

type map2Input struct {
	id       string
	time_col string
}

type map2Expected struct {
	id       uint64
	time_col time.Time
}

type map3Input struct {
	id         string
	string_col string
	binary_col string
	int_col    string
	uint_col   string
	float_col  string
	bool_col   string
	any_col    string
}

type map3Expected struct {
	id         uint64
	string_col string
	binary_col []byte
	int_col    int64
	uint_col   uint64
	float_col  float64
	bool_col   bool
	any_col    interface{}
}

var (
	SPECIAL_CHAR                      = `"!@#$%^&*()_+{}:<>?~-=[]',./;"`
	SPECIAL_CHAR_EXPECTED_OUTPUT      = `!@#$%^&*()_+{}:<>?~-=[]',./;`
	SPECIAL_CHAR_EXPECTED_BYTE_OUTPUT = []byte{0x21, 0x40, 0x23, 0x24, 0x25, 0x5e, 0x26, 0x2a, 0x28, 0x29, 0x5f, 0x2b, 0x7b, 0x7d, 0x3a, 0x3c, 0x3e, 0x3f, 0x7e, 0x2d, 0x3d, 0x5b, 0x5d, 0x27, 0x2c, 0x2e, 0x2f, 0x3b}
	LARGE_STRING_EXPECTED_OUTPUT      string
	LARGE_BYTE_EXPECTED_OUTPUT        []byte
	LARGE_STRING                      string
	LARGE_BYTE                        string
)

func init() {

	LARGE_STRING_EXPECTED_OUTPUT = ""
	LARGE_BYTE_EXPECTED_OUTPUT = make([]byte, 32768)
	LARGE_BYTE = ""
	for i := 0; i < 32768; i++ {
		LARGE_STRING_EXPECTED_OUTPUT += "a"
		LARGE_BYTE_EXPECTED_OUTPUT[i] = 0x30
		LARGE_BYTE += "0"
	}

	LARGE_STRING = fmt.Sprintf(`"%s"`, LARGE_STRING_EXPECTED_OUTPUT)
	LARGE_BYTE = fmt.Sprintf(`"%s"`, LARGE_BYTE)
}

func MapTransformTest2(t *testing.T, testSourceSchema *schemer.Schema, transformer *schemer.Transformer, input map2Input, expected map2Expected) {

	jsonInput := fmt.Sprintf(`
	{
		"id":           %s,
		"map_col": {
			"time_col": %s
		}
	}`, input.id, input.time_col)
	var rawData map[string]interface{}
	err := json.Unmarshal([]byte(jsonInput), &rawData)
	if err != nil {
		t.Fatal(err)
	}

	source := testSourceSchema.Normalize(rawData)

	if err != nil {
		t.Fatal(err)
	}

	returnedValue, err := transformer.Transform(nil, source)
	if !assert.Nil(t, err) {
		t.Fatalf("transform failed: %v", err)
	}

	if !assert.Len(t, returnedValue, 1) {
		t.Fatalf("return length not match")
	}

	result := returnedValue[0]
	AssertMap2Result(t, result, expected)
}

func MapTransformTest3(t *testing.T, testSourceSchema *schemer.Schema, transformer *schemer.Transformer, input map3Input, expected map3Expected) {

	jsonInput := fmt.Sprintf(`
	{
		"id":             %s,
		"map_col": {
			"string_col": %s,
			"binary_col": %s,
			"int_col":    %s,
			"uint_col":   %s,
			"float_col":  %s,
			"bool_col":   %s,
			"any_col":    %s
		}
	}`, input.id, input.string_col, input.binary_col, input.int_col, input.uint_col, input.float_col, input.bool_col, input.any_col)
	var rawData map[string]interface{}
	err := json.Unmarshal([]byte(jsonInput), &rawData)
	if err != nil {
		t.Fatal(err)
	}

	source := testSourceSchema.Normalize(rawData)

	returnedValue, err := transformer.Transform(nil, source)
	if !assert.Nil(t, err) {
		t.Fatalf("transform failed: %v", err)
	}

	if !assert.Len(t, returnedValue, 1) {
		t.Fatalf("return length not match")
	}

	result := returnedValue[0]
	AssertMap3Result(t, result, expected)
}

func AssertMap2Result(t *testing.T, result map[string]interface{}, expected map2Expected) {

	assert.Equal(t, expected.id, result["id"])
	map_col := result["map_col"].(map[string]interface{})
	assert.Equal(t, expected.time_col.UTC(), map_col["time_col"].(time.Time).UTC())
}

func AssertMap3Result(t *testing.T, result map[string]interface{}, expected map3Expected) {

	assert.Equal(t, expected.id, result["id"])
	map_col := result["map_col"].(map[string]interface{})
	assert.Equal(t, expected.string_col, map_col["string_col"])
	assert.Equal(t, expected.binary_col, map_col["binary_col"])
	assert.Equal(t, expected.int_col, map_col["int_col"])
	assert.Equal(t, expected.uint_col, map_col["uint_col"])
	assert.Equal(t, expected.float_col, map_col["float_col"])
	assert.Equal(t, expected.bool_col, map_col["bool_col"])
	assert.Equal(t, expected.any_col, map_col["any_col"])
}

func TestMapEmpty(t *testing.T) {

	testSourceSchema := schemer.NewSchema()
	err := schemer.UnmarshalJSON([]byte(mapSchema1), testSourceSchema)
	assert.Error(t, err)

	testDestSchema := schemer.NewSchema()
	err = schemer.UnmarshalJSON([]byte(mapSchema1), testDestSchema)
	assert.Error(t, err)
}

func TestMapSuccessTransform2(t *testing.T) {

	testSourceSchema := schemer.NewSchema()
	err := schemer.UnmarshalJSON([]byte(mapSchema2), testSourceSchema)
	if err != nil {
		t.Error(err)
	}

	// Using the same schema for destination
	testDestSchema := schemer.NewSchema()
	err = schemer.UnmarshalJSON([]byte(mapSchema2), testDestSchema)
	if err != nil {
		t.Error(err)
	}

	// Create transformer
	transformer := schemer.NewTransformer(testSourceSchema, testDestSchema)
	err = transformer.SetScript(`return source`)
	if err != nil {
		t.Error(err)
	}

	maptest1 := map2Input{`1`, `"2024-08-06T15:02:00+08:00"`}
	maptest1Expected := map2Expected{1, time.Date(2024, 8, 6, 7, 2, 0, 0, time.UTC)}
	MapTransformTest2(t, testSourceSchema, transformer, maptest1, maptest1Expected)

	maptest2 := map2Input{`2`, `"2024-08-06T15:02:00Z"`}
	maptest2Expected := map2Expected{2, time.Date(2024, 8, 6, 15, 2, 0, 0, time.UTC)}
	MapTransformTest2(t, testSourceSchema, transformer, maptest2, maptest2Expected)
}

func TestMapTransformErrorHandle2(t *testing.T) {

	testSourceSchema := schemer.NewSchema()
	err := schemer.UnmarshalJSON([]byte(mapSchema2), testSourceSchema)
	if err != nil {
		t.Error(err)
	}

	// Using the same schema for destination
	testDestSchema := schemer.NewSchema()
	err = schemer.UnmarshalJSON([]byte(mapSchema2), testDestSchema)
	if err != nil {
		t.Error(err)
	}

	// Create transformer
	transformer := schemer.NewTransformer(testSourceSchema, testDestSchema)
	err = transformer.SetScript(`return source`)
	if err != nil {
		t.Error(err)
	}

	maptest1 := map2Input{`1`, `"2024-08-06T15:02:00.1234567890Z"`}
	maptest1Expected := map2Expected{1, time.Date(2024, 8, 6, 15, 2, 0, 123000000, time.UTC)}
	MapTransformTest2(t, testSourceSchema, transformer, maptest1, maptest1Expected)
}

func TestMapSuccessTransform3(t *testing.T) {

	testSourceSchema := schemer.NewSchema()
	err := schemer.UnmarshalJSON([]byte(mapSchema3), testSourceSchema)
	if err != nil {
		t.Error(err)
	}

	// Using the same schema for destination
	testDestSchema := schemer.NewSchema()
	err = schemer.UnmarshalJSON([]byte(mapSchema3), testDestSchema)
	if err != nil {
		t.Error(err)
	}

	// Create transformer
	transformer := schemer.NewTransformer(testSourceSchema, testDestSchema)
	err = transformer.SetScript(`return source`)
	if err != nil {
		t.Error(err)
	}

	maptest1 := map3Input{`1`, `""`, `""`, `5`, `5`, `5`, `0`, `""`}
	maptest1Expected := map3Expected{1, "", []byte{}, 5, 5, 5, false, ""}
	MapTransformTest3(t, testSourceSchema, transformer, maptest1, maptest1Expected)

	maptest2 := map3Input{`2`, `" "`, `" "`, `0`, `0`, `1.23`, `1`, `" "`}
	maptest2Expected := map3Expected{2, " ", []byte{0x20}, 0, 0, 1.23, true, " "}
	MapTransformTest3(t, testSourceSchema, transformer, maptest2, maptest2Expected)

	maptest3 := map3Input{`3`, `"abc"`, LARGE_BYTE, `-1`, `5`, `-1.23`, `"false"`, `"abc"`}
	maptest3Expected := map3Expected{3, "abc", LARGE_BYTE_EXPECTED_OUTPUT, -1, 5, -1.23, false, "abc"}
	MapTransformTest3(t, testSourceSchema, transformer, maptest3, maptest3Expected)

	maptest4 := map3Input{`4`, `"中文"`, `"0"`, `5`, `0`, `-1.234567111111111`, `"true"`, `"中文"`}
	maptest4Expected := map3Expected{4, "中文", []byte{0x30}, 5, 0, -1.234567111111111, true, "中文"}
	MapTransformTest3(t, testSourceSchema, transformer, maptest4, maptest4Expected)

	maptest5 := map3Input{`5`, SPECIAL_CHAR, `"001"`, `0`, `5`, `1.234567111111111`, `"True"`, SPECIAL_CHAR}
	maptest5Expected := map3Expected{5, SPECIAL_CHAR_EXPECTED_OUTPUT, []byte{0x30, 0x30, 0x31}, 0, 5, 1.234567111111111, true, SPECIAL_CHAR_EXPECTED_OUTPUT}
	MapTransformTest3(t, testSourceSchema, transformer, maptest5, maptest5Expected)

	maptest6 := map3Input{`6`, LARGE_STRING, `""`, `-1`, `0`, `1.7976931348623157e+308`, `"False"`, LARGE_STRING}
	maptest6Expected := map3Expected{6, LARGE_STRING_EXPECTED_OUTPUT, []byte{}, -1, 0, 1.7976931348623157e+308, false, LARGE_STRING_EXPECTED_OUTPUT}
	MapTransformTest3(t, testSourceSchema, transformer, maptest6, maptest6Expected)

	maptest7 := map3Input{`7`, `""`, `" "`, `5`, `5`, `-1.7976931348623157e+308`, `"T"`, `5`}
	maptest7Expected := map3Expected{7, "", []byte{0x20}, 5, 5, -1.7976931348623157e+308, true, float64(5)}
	MapTransformTest3(t, testSourceSchema, transformer, maptest7, maptest7Expected)

	maptest8 := map3Input{`8`, `" "`, `"0"`, `0`, `0`, `-0`, `"F"`, `[]`}
	maptest8Expected := map3Expected{8, " ", []byte{0x30}, 0, 0, 0, false, []interface{}{}}
	MapTransformTest3(t, testSourceSchema, transformer, maptest8, maptest8Expected)

	maptest9 := map3Input{`9`, `"abc"`, `"001"`, `-1`, `5`, `5`, `"t"`, `{}`}
	maptest9Expected := map3Expected{9, "abc", []byte{0x30, 0x30, 0x31}, -1, 5, 5, true, map[string]interface{}{}}
	MapTransformTest3(t, testSourceSchema, transformer, maptest9, maptest9Expected)

	maptest10 := map3Input{`10`, `"中文"`, `""`, `5`, `0`, `1.23`, `"f"`, `true`}
	maptest10Expected := map3Expected{10, "中文", []byte{}, 5, 0, 1.23, false, true}
	MapTransformTest3(t, testSourceSchema, transformer, maptest10, maptest10Expected)

	maptest11 := map3Input{`11`, SPECIAL_CHAR, `" "`, `0`, `5`, `-1.23`, `"0"`, `null`}
	maptest11Expected := map3Expected{11, SPECIAL_CHAR_EXPECTED_OUTPUT, []byte{0x20}, 0, 5, -1.23, false, nil}
	MapTransformTest3(t, testSourceSchema, transformer, maptest11, maptest11Expected)

	maptest12 := map3Input{`12`, LARGE_STRING, LARGE_BYTE, `-1`, `0`, `-1.234567111111111`, `"1"`, `""`}
	maptest12Expected := map3Expected{12, LARGE_STRING_EXPECTED_OUTPUT, LARGE_BYTE_EXPECTED_OUTPUT, -1, 0, -1.234567111111111, true, ""}
	MapTransformTest3(t, testSourceSchema, transformer, maptest12, maptest12Expected)
}

func TestMapTransformErrorHandle3(t *testing.T) {

	testSourceSchema := schemer.NewSchema()
	err := schemer.UnmarshalJSON([]byte(mapSchema3), testSourceSchema)
	if err != nil {
		t.Error(err)
	}

	// Using the same schema for destination
	testDestSchema := schemer.NewSchema()
	err = schemer.UnmarshalJSON([]byte(mapSchema3), testDestSchema)
	if err != nil {
		t.Error(err)
	}

	// Create transformer
	transformer := schemer.NewTransformer(testSourceSchema, testDestSchema)
	err = transformer.SetScript(`return source`)
	if err != nil {
		t.Error(err)
	}

	maptest1 := map3Input{`1`, `5`, `"abc"`, `""`, `""`, `""`, `""`, `""`}
	maptest1Expected := map3Expected{1, "5", []byte{0x61, 0x62, 0x63}, 0, 0, 0, false, ""}
	MapTransformTest3(t, testSourceSchema, transformer, maptest1, maptest1Expected)

	maptest2 := map3Input{`2`, `5`, `"中文"`, `" "`, `" "`, `" "`, `" "`, `""`}
	maptest2Expected := map3Expected{2, "5", []byte{0xe4, 0xb8, 0xad, 0xe6, 0x96, 0x87}, 0, 0, 0, false, ""}
	MapTransformTest3(t, testSourceSchema, transformer, maptest2, maptest2Expected)

	maptest3 := map3Input{`3`, `5`, SPECIAL_CHAR, `"abc"`, `"abc"`, `"abc"`, `"abc"`, `""`}
	maptest3Expected := map3Expected{3, "5", SPECIAL_CHAR_EXPECTED_BYTE_OUTPUT, 0, 0, 0, false, ""}
	MapTransformTest3(t, testSourceSchema, transformer, maptest3, maptest3Expected)

	maptest4 := map3Input{`4`, `5`, `5`, `"中文"`, `"中文"`, `"中文"`, `"中文"`, `""`}
	maptest4Expected := map3Expected{4, "5", []byte{}, 0, 0, 0, false, ""}
	MapTransformTest3(t, testSourceSchema, transformer, maptest4, maptest4Expected)

	maptest5 := map3Input{`5`, `5`, `"10102"`, SPECIAL_CHAR, SPECIAL_CHAR, SPECIAL_CHAR, SPECIAL_CHAR, `""`}
	maptest5Expected := map3Expected{5, "5", []byte{0x31, 0x30, 0x31, 0x30, 0x32}, 0, 0, 0, false, ""}
	MapTransformTest3(t, testSourceSchema, transformer, maptest5, maptest5Expected)

	maptest6 := map3Input{`6`, `5`, `101`, LARGE_STRING, LARGE_STRING, LARGE_STRING, LARGE_STRING, `""`}
	maptest6Expected := map3Expected{6, "5", []byte{}, 0, 0, 0, false, ""}
	MapTransformTest3(t, testSourceSchema, transformer, maptest6, maptest6Expected)

	// 以下大於 int64 的數值會變成最小負值，uint64(-1) 會變成 最大uint，小於float64最小值的數值會去掉小數點，bool小於0會變成false，大於0會變成true
	maptest7 := map3Input{`7`, `5`, `"abc"`, `9223372036854775808`, `-1`, `1.0000000000000001`, `5`, `""`}
	maptest7Expected := map3Expected{7, "5", []byte{0x61, 0x62, 0x63}, -9223372036854775808, 18446744073709551615, 1, true, ""}
	MapTransformTest3(t, testSourceSchema, transformer, maptest7, maptest7Expected)

	// 小於int64最小值會被鎖定在最小值，uint大於一定數會變固定值
	maptest8 := map3Input{`8`, `5`, `"中文"`, `-9223372036854775809`, `18446744073709551615`, `""`, `""`, `""`}
	maptest8Expected := map3Expected{8, "5", []byte{0xe4, 0xb8, 0xad, 0xe6, 0x96, 0x87}, 9223372036854775807, 18446744073709551615, 0, false, ""}
	MapTransformTest3(t, testSourceSchema, transformer, maptest8, maptest8Expected)

	// int與uint帶浮點數會被去掉小數點
	maptest9 := map3Input{`9`, `5`, SPECIAL_CHAR, `1.23`, `1.23`, `" "`, `" "`, `""`}
	maptest9Expected := map3Expected{9, "5", SPECIAL_CHAR_EXPECTED_BYTE_OUTPUT, 1, 1, 0, false, ""}
	MapTransformTest3(t, testSourceSchema, transformer, maptest9, maptest9Expected)

	// 測試重複key與缺失key情形
	jsonInput := fmt.Sprintf(`
	{
		"id":             %s,
		"map_col": {
			"string_col": %s,
			"string_col": %s
		}
	}`, `13`, `"abc1"`, `"abc2"`)
	var rawData map[string]interface{}
	err = json.Unmarshal([]byte(jsonInput), &rawData)
	if err != nil {
		t.Error(err)
		return
	}

	source := testSourceSchema.Normalize(rawData)
	returnedValue, err := transformer.Transform(nil, source)
	if !assert.Nil(t, err) {
		t.Fatalf("transform failed: %v", err)
	}

	if !assert.Len(t, returnedValue, 1) {
		t.Fatalf("return length not match")
	}

	result := returnedValue[0]
	assert.Equal(t, uint64(13), result["id"])
	map_col := result["map_col"].(map[string]interface{})
	assert.Equal(t, "abc2", map_col["string_col"])
	assert.Equal(t, nil, map_col["binary_col"])
}
