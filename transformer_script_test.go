package schemer

import (
	"testing"

	"github.com/dop251/goja"
	"github.com/stretchr/testify/assert"
)

func TestTransformerScript_ScanStruct(t *testing.T) {

	script := `function test() {
		let input = {
			"has_value": "value",
			"undefined_value": undefined,
			"null_value": null,
			"arr": [
				1,
				"string",
				undefined,
				null
			],
			"obj": {
				"key1": "value1",
				"key2": "value2",
				"key3": undefined
			}
		};
		scanStruct(input);
		return input;
}`

	// Initialize VM and inject functions
	vm := goja.New()
	ts := NewTransformerScript()
	ts.injectFuncs(vm)
	vm.RunString(script)

	// Run test function
	testFunc, _ := goja.AssertFunction(vm.Get("test"))
	res, err := testFunc(goja.Undefined())
	if !assert.Nil(t, err) {
		t.Error(err)
		return
	}

	assert.False(t, goja.IsUndefined(res))

	obj := res.Export().(map[string]interface{})

	assert.Equal(t, "value", obj["has_value"])
	assert.NotContains(t, obj, "undefined_value")
	assert.Contains(t, obj, "null_value")

	// arr
	arr := obj["arr"].([]interface{})
	assert.Equal(t, arr[0], int64(1))
	assert.Equal(t, arr[1], "string")
	assert.Equal(t, arr[2], nil)
	assert.Equal(t, arr[3], nil)

	// obj
	subObj := obj["obj"].(map[string]interface{})
	assert.Equal(t, subObj["key1"], "value1")
	assert.Equal(t, subObj["key2"], "value2")
	assert.NotContains(t, subObj, "key3")
}
