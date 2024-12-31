package schemer_test

import (
	"encoding/json"
	"fmt"
	"math"
	"testing"

	"github.com/BrobridgeOrg/schemer"
	"github.com/stretchr/testify/assert"
)

type normalSchemaInput struct {
	id         string
	string_col string
	binary_col string
	int_col    string
	uint_col   string
	float_col  string
	bool_col   string
	any_col    string
}

type normalSchemaExpected struct {
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
	NORMAL_SCHEMA_SPECIAL_CHAR                      = `"!@#$%^&*()_+{}:<>?~-=[]',./;"`
	NORMAL_SCHEMA_SPECIAL_CHAR_EXPECTED_OUTPUT      = `!@#$%^&*()_+{}:<>?~-=[]',./;`
	NORMAL_SCHEMA_SPECIAL_CHAR_EXPECTED_BYTE_OUTPUT = []byte{0x21, 0x40, 0x23, 0x24, 0x25, 0x5e, 0x26, 0x2a, 0x28, 0x29, 0x5f, 0x2b, 0x7b, 0x7d, 0x3a, 0x3c, 0x3e, 0x3f, 0x7e, 0x2d, 0x3d, 0x5b, 0x5d, 0x27, 0x2c, 0x2e, 0x2f, 0x3b}
	NORMAL_SCHEMA_LARGE_STRING_EXPECTED_OUTPUT      string
	NORMAL_SCHEMA_LARGE_BYTE_EXPECTED_OUTPUT        []byte
	NORMAL_SCHEMA_LARGE_STRING                      string
	NORMAL_SCHEMA_LARGE_BYTE                        string
)

var normalSchema = `{
	"id":{
	   "type":"uint"
	},
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
}`

func init() {

	NORMAL_SCHEMA_LARGE_STRING_EXPECTED_OUTPUT = ""
	NORMAL_SCHEMA_LARGE_BYTE_EXPECTED_OUTPUT = make([]byte, 32768)
	NORMAL_SCHEMA_LARGE_BYTE = ""
	for i := 0; i < 32768; i++ {
		NORMAL_SCHEMA_LARGE_STRING += "a"
		NORMAL_SCHEMA_LARGE_STRING_EXPECTED_OUTPUT += "a"
		NORMAL_SCHEMA_LARGE_BYTE_EXPECTED_OUTPUT[i] = 0x30
		NORMAL_SCHEMA_LARGE_BYTE += "0"
	}
	NORMAL_SCHEMA_LARGE_STRING = fmt.Sprintf(`"%s"`, NORMAL_SCHEMA_LARGE_STRING)
	NORMAL_SCHEMA_LARGE_BYTE = fmt.Sprintf(`"%s"`, NORMAL_SCHEMA_LARGE_BYTE)
}

func NormalizeNormalSchema(s *schemer.Schema, input normalSchemaInput) (map[string]interface{}, error) {

	jsonInput := fmt.Sprintf(`
	{
		"id":         %s,
		"string_col": %s,
		"binary_col": %s,
		"int_col":    %s,
		"uint_col":   %s,
		"float_col":  %s,
		"bool_col":   %s,
		"any_col":    %s
	}`, input.id, input.string_col, input.binary_col, input.int_col, input.uint_col, input.float_col, input.bool_col, input.any_col)
	var rawData map[string]interface{}
	err := json.Unmarshal([]byte(jsonInput), &rawData)
	if err != nil {
		return nil, err
	}
	return s.Normalize(rawData), nil
}

func NormalSchemaTransformTest(t *testing.T, testSourceSchema *schemer.Schema, transformer *schemer.Transformer, input normalSchemaInput, expected normalSchemaExpected) {

	source, err := NormalizeNormalSchema(testSourceSchema, input)
	if err != nil {
		t.Fatal(err)
	}

	returnedValue, err := transformer.Transform(nil, source)
	if !assert.Nil(t, err) {
		t.Fatal(err)
	}

	if !assert.Len(t, returnedValue, 1) {
		t.Fatal(err)
	}

	result := returnedValue[0]
	AssertNormalSchemaResult(t, result, expected)
}

func AssertNormalSchemaResult(t *testing.T, result map[string]interface{}, expected normalSchemaExpected) {

	assert.Equal(t, expected.id, result["id"])
	assert.Equal(t, expected.string_col, result["string_col"])
	assert.Equal(t, expected.binary_col, result["binary_col"])
	assert.Equal(t, expected.int_col, result["int_col"])
	assert.Equal(t, expected.uint_col, result["uint_col"])
	assert.Equal(t, expected.float_col, result["float_col"])
	assert.Equal(t, expected.bool_col, result["bool_col"])
	assert.Equal(t, expected.any_col, result["any_col"])
}

func TestNormalSchema(t *testing.T) {

	testSourceSchema := schemer.NewSchema()
	err := schemer.UnmarshalJSON([]byte(normalSchema), testSourceSchema)
	if err != nil {
		t.Error(err)
	}

	testDestSchema := schemer.NewSchema()
	err = schemer.UnmarshalJSON([]byte(normalSchema), testDestSchema)
	if err != nil {
		t.Error(err)
	}

	transformer := schemer.NewTransformer(testSourceSchema, testDestSchema, schemer.WithRuntime(jsRuntime))
	err = transformer.SetScript(`return source`)
	if err != nil {
		t.Error(err)
	}

	mainSuccessInput1 := normalSchemaInput{`1`, `""`, `""`, `5`, `5`, `5`, `0`, `""`}
	mainSuccessExpected1 := normalSchemaExpected{1, "", []byte{}, 5, 5, 5, false, ""}
	NormalSchemaTransformTest(t, testSourceSchema, transformer, mainSuccessInput1, mainSuccessExpected1)

	mainSuccessInput2 := normalSchemaInput{`2`, `" "`, `" "`, `0`, `0`, `1.23`, `1`, `" "`}
	mainSuccessExpected2 := normalSchemaExpected{2, " ", []byte{0x20}, 0, 0, 1.23, true, " "}
	NormalSchemaTransformTest(t, testSourceSchema, transformer, mainSuccessInput2, mainSuccessExpected2)

	mainSuccessInput3 := normalSchemaInput{`3`, `"abc"`, NORMAL_SCHEMA_LARGE_BYTE, `-1`, `5`, `-1.23`, `"false"`, `"abc"`}
	mainSuccessExpected3 := normalSchemaExpected{3, "abc", NORMAL_SCHEMA_LARGE_BYTE_EXPECTED_OUTPUT, -1, 5, -1.23, false, "abc"}
	NormalSchemaTransformTest(t, testSourceSchema, transformer, mainSuccessInput3, mainSuccessExpected3)

	mainSuccessInput4 := normalSchemaInput{`4`, `"中文"`, `"0"`, `5`, `0`, `-1.234567111111111`, `"true"`, `"中文"`}
	mainSuccessExpected4 := normalSchemaExpected{4, "中文", []byte{0x30}, 5, 0, -1.234567111111111, true, "中文"}
	NormalSchemaTransformTest(t, testSourceSchema, transformer, mainSuccessInput4, mainSuccessExpected4)

	mainSuccessInput5 := normalSchemaInput{`5`, NORMAL_SCHEMA_SPECIAL_CHAR, `"001"`, `0`, `5`, `1.234567111111111`, `"True"`, NORMAL_SCHEMA_SPECIAL_CHAR}
	mainSuccessExpected5 := normalSchemaExpected{5, NORMAL_SCHEMA_SPECIAL_CHAR_EXPECTED_OUTPUT, []byte{0x30, 0x30, 0x31}, 0, 5, 1.234567111111111, true, NORMAL_SCHEMA_SPECIAL_CHAR_EXPECTED_OUTPUT}
	NormalSchemaTransformTest(t, testSourceSchema, transformer, mainSuccessInput5, mainSuccessExpected5)

	mainSuccessInput6 := normalSchemaInput{`6`, NORMAL_SCHEMA_LARGE_STRING, `""`, `-1`, `0`, `1.7976931348623157e+308`, `"False"`, `""`}
	mainSuccessExpected6 := normalSchemaExpected{6, NORMAL_SCHEMA_LARGE_STRING_EXPECTED_OUTPUT, []byte{}, -1, 0, math.MaxFloat64, false, ""}
	NormalSchemaTransformTest(t, testSourceSchema, transformer, mainSuccessInput6, mainSuccessExpected6)

	mainSuccessInput7 := normalSchemaInput{`7`, `""`, `" "`, `5`, `5`, `-1.7976931348623157e+308`, `"T"`, `5`}
	mainSuccessExpected7 := normalSchemaExpected{7, "", []byte{0x20}, 5, 5, -math.MaxFloat64, true, int64(5)}
	NormalSchemaTransformTest(t, testSourceSchema, transformer, mainSuccessInput7, mainSuccessExpected7)

	mainSuccessInput8 := normalSchemaInput{`8`, `" "`, `"0"`, `0`, `0`, `-0`, `"F"`, `[]`}
	mainSuccessExpected8 := normalSchemaExpected{8, " ", []byte{0x30}, 0, 0, -0, false, []interface{}{}}
	NormalSchemaTransformTest(t, testSourceSchema, transformer, mainSuccessInput8, mainSuccessExpected8)

	mainSuccessInput9 := normalSchemaInput{`9`, `"abc"`, `"001"`, `-1`, `5`, `5`, `"t"`, `{}`}
	mainSuccessExpected9 := normalSchemaExpected{9, "abc", []byte{0x30, 0x30, 0x31}, -1, 5, 5, true, map[string]interface{}{}}
	NormalSchemaTransformTest(t, testSourceSchema, transformer, mainSuccessInput9, mainSuccessExpected9)

	mainSuccessInput10 := normalSchemaInput{`10`, `"中文"`, `""`, `5`, `0`, `1.23`, `"f"`, `true`}
	mainSuccessExpected10 := normalSchemaExpected{10, "中文", []byte{}, 5, 0, 1.23, false, true}
	NormalSchemaTransformTest(t, testSourceSchema, transformer, mainSuccessInput10, mainSuccessExpected10)

	mainSuccessInput11 := normalSchemaInput{`11`, NORMAL_SCHEMA_SPECIAL_CHAR, `" "`, `0`, `5`, `-1.23`, `"0"`, `null`}
	mainSuccessExpected11 := normalSchemaExpected{11, NORMAL_SCHEMA_SPECIAL_CHAR_EXPECTED_OUTPUT, []byte{0x20}, 0, 5, -1.23, false, nil}
	NormalSchemaTransformTest(t, testSourceSchema, transformer, mainSuccessInput11, mainSuccessExpected11)

	mainSuccessInput12 := normalSchemaInput{`12`, NORMAL_SCHEMA_LARGE_STRING, NORMAL_SCHEMA_LARGE_BYTE, `-1`, `0`, `-1.234567111111111`, `"1"`, `""`}
	mainSuccessExpected12 := normalSchemaExpected{12, NORMAL_SCHEMA_LARGE_STRING_EXPECTED_OUTPUT, NORMAL_SCHEMA_LARGE_BYTE_EXPECTED_OUTPUT, -1, 0, -1.234567111111111, true, ""}
	NormalSchemaTransformTest(t, testSourceSchema, transformer, mainSuccessInput12, mainSuccessExpected12)
}

func TestNormalSchemaErrorHandle(t *testing.T) {

	testSourceSchema := schemer.NewSchema()
	err := schemer.UnmarshalJSON([]byte(normalSchema), testSourceSchema)
	if err != nil {
		t.Error(err)
	}

	testDestSchema := schemer.NewSchema()
	err = schemer.UnmarshalJSON([]byte(normalSchema), testDestSchema)
	if err != nil {
		t.Error(err)
	}

	transformer := schemer.NewTransformer(testSourceSchema, testDestSchema, schemer.WithRuntime(jsRuntime))
	err = transformer.SetScript(`return source`)
	if err != nil {
		t.Error(err)
	}

	extensionOneInput1 := normalSchemaInput{`1`, `5`, `"abc"`, `""`, `""`, `""`, `""`, `""`}
	extensionOneExpected1 := normalSchemaExpected{1, "5", []byte{0x61, 0x62, 0x63}, 0, 0, 0, false, ""}
	NormalSchemaTransformTest(t, testSourceSchema, transformer, extensionOneInput1, extensionOneExpected1)

	extensionOneInput2 := normalSchemaInput{`2`, `5`, `"中文"`, `" "`, `" "`, `" "`, `" "`, `""`}
	extensionOneExpected2 := normalSchemaExpected{2, "5", []byte{0xe4, 0xb8, 0xad, 0xe6, 0x96, 0x87}, 0, 0, 0, false, ""}
	NormalSchemaTransformTest(t, testSourceSchema, transformer, extensionOneInput2, extensionOneExpected2)

	extensionOneInput3 := normalSchemaInput{`3`, `5`, NORMAL_SCHEMA_SPECIAL_CHAR, `"abc"`, `"abc"`, `"abc"`, `"abc"`, `""`}
	extensionOneExpected3 := normalSchemaExpected{3, "5", NORMAL_SCHEMA_SPECIAL_CHAR_EXPECTED_BYTE_OUTPUT, 0, 0, 0, false, ""}
	NormalSchemaTransformTest(t, testSourceSchema, transformer, extensionOneInput3, extensionOneExpected3)

	extensionOneInput4 := normalSchemaInput{`4`, `5`, `5`, `"中文"`, `"中文"`, `"中文"`, `"中文"`, `""`}
	extensionOneExpected4 := normalSchemaExpected{4, "5", []byte{}, 0, 0, 0, false, ""}
	NormalSchemaTransformTest(t, testSourceSchema, transformer, extensionOneInput4, extensionOneExpected4)

	extensionOneInput5 := normalSchemaInput{`5`, `5`, `"10102"`, NORMAL_SCHEMA_SPECIAL_CHAR, NORMAL_SCHEMA_SPECIAL_CHAR, NORMAL_SCHEMA_SPECIAL_CHAR, NORMAL_SCHEMA_SPECIAL_CHAR, `""`}
	extensionOneExpected5 := normalSchemaExpected{5, "5", []byte{0x31, 0x30, 0x31, 0x30, 0x32}, 0, 0, 0, false, ""}
	NormalSchemaTransformTest(t, testSourceSchema, transformer, extensionOneInput5, extensionOneExpected5)

	extensionOneInput6 := normalSchemaInput{`6`, `5`, `101`, NORMAL_SCHEMA_LARGE_STRING, NORMAL_SCHEMA_LARGE_STRING, NORMAL_SCHEMA_LARGE_STRING, NORMAL_SCHEMA_LARGE_STRING, `""`}
	extensionOneExpected6 := normalSchemaExpected{6, "5", []byte{}, 0, 0, 0, false, ""}
	NormalSchemaTransformTest(t, testSourceSchema, transformer, extensionOneInput6, extensionOneExpected6)

	extensionOneInput7 := normalSchemaInput{`7`, `5`, `"abc"`, `9223372036854775808`, `-1`, `1.0000000000000001`, `5`, `""`}
	extensionOneExpected7 := normalSchemaExpected{7, "5", []byte{0x61, 0x62, 0x63}, math.MinInt64, math.MaxUint64, 1, true, ""}
	NormalSchemaTransformTest(t, testSourceSchema, transformer, extensionOneInput7, extensionOneExpected7)

	extensionOneInput8 := normalSchemaInput{`8`, `5`, `"中文"`, `-9223372036854775809`, `18446744073709551616`, `""`, `""`, `""`}
	extensionOneExpected8 := normalSchemaExpected{8, "5", []byte{0xe4, 0xb8, 0xad, 0xe6, 0x96, 0x87}, math.MaxInt64, 0, 0, false, ""}
	NormalSchemaTransformTest(t, testSourceSchema, transformer, extensionOneInput8, extensionOneExpected8)

	extensionOneInput9 := normalSchemaInput{`9`, `5`, NORMAL_SCHEMA_SPECIAL_CHAR, `1.23`, `1.23`, `" "`, `" "`, `""`}
	extensionOneExpected9 := normalSchemaExpected{9, "5", NORMAL_SCHEMA_SPECIAL_CHAR_EXPECTED_BYTE_OUTPUT, 1, 1, 0, false, ""}
	NormalSchemaTransformTest(t, testSourceSchema, transformer, extensionOneInput9, extensionOneExpected9)
}
