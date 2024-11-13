package tests

import (
	"regexp"
	"strings"
	"testing"

	"github.com/BrobridgeOrg/schemer"
	"github.com/stretchr/testify/assert"
)

var arraySchema2 = `{
    "array_null":{
        "type":"array",
        "subtype":""
    },
    "array_space":{
        "type":"array",
        "subtype":" "
    },
    "array_abc":{
        "type":"array",
        "subtype":"abc"
    },
    "array_chinese":{
        "type":"array",
        "subtype":"中文"
    },
    "array_special":{
        "type":"array",
        "subtype":"!@#$%^&*()_+{}:<>?~-=[]\\;',./"
    },
    "array_maxLen":{
        "type":"array",
        "subtype":"[max_len_str()]"
    },
    "array_Ignore":{
        "type":"array"
    }
}`

var arraySchema3 = `{
    "array_null":{
        "type":"array",
        "subtype":""
    }
}`

var arraySchema4 = `{
    "array_space":{
        "type":"array",
        "subtype":" "
    }
}`

var arraySchema5 = `{
    "array_abc":{
        "type":"array",
        "subtype":"abc"
    }
}`

var arraySchema6 = `{ 
    "array_chinese":{
        "type":"array",
        "subtype":"中文"
    }
}`

var arraySchema7 = `{
    "array_special":{
        "type":"array",
        "subtype":"!@#$%^&*()_+{}:<>?~-=[]\\;',./" 
    }
}`
var arraySchema8 = `{
    "array_maxLen":{
        "type":"array",
        "subtype":"[max_len_str()]"
    }
}`

var arraySchema9 = `{
    "array_Ignore":{
        "type":"array"
    }
}`

func TestWorngSubtypeWithString(t *testing.T) {

	testSourceSchema := schemer.NewSchema()

	arraySchema2 = ReplaceMaxLenStr(arraySchema2)
	err := schemer.UnmarshalJSON([]byte(arraySchema2), testSourceSchema)
	assert.Error(t, err)

	err = schemer.UnmarshalJSON([]byte(arraySchema3), testSourceSchema)
	assert.Error(t, err)

	err = schemer.UnmarshalJSON([]byte(arraySchema4), testSourceSchema)
	assert.Error(t, err)

	err = schemer.UnmarshalJSON([]byte(arraySchema5), testSourceSchema)
	assert.Error(t, err)

	err = schemer.UnmarshalJSON([]byte(arraySchema6), testSourceSchema)
	assert.Error(t, err)

	err = schemer.UnmarshalJSON([]byte(arraySchema7), testSourceSchema)
	assert.Error(t, err)

	arraySchema8 = ReplaceMaxLenStr(arraySchema8)
	err = schemer.UnmarshalJSON([]byte(arraySchema8), testSourceSchema)
	assert.Error(t, err)

	err = schemer.UnmarshalJSON([]byte(arraySchema9), testSourceSchema)
	assert.Error(t, err)
}

func ReplaceMaxLenStr(s string) string {

	reMaxLenStr := regexp.MustCompile(`\[max_len_str\(\)\]`)
	if reMaxLenStr.MatchString(s) {
		longString := strings.Repeat("a", 32768)
		s = reMaxLenStr.ReplaceAllString(s, longString)
	}
	return s
}
