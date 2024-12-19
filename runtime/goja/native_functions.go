package goja_runtime

import (
	"reflect"
	"strings"

	"github.com/dop251/goja"
)

func (r *Runtime) initNativeFunctions() error {

	r.vm.Set("prepareRefs", func(call goja.FunctionCall) goja.Value {
		input := call.Argument(0).Export().(map[string]interface{})
		result := r.nativePrepareRefs(input)
		return r.vm.ToValue(result)
	})

	r.vm.Set("scanStruct", func(call goja.FunctionCall) goja.Value {
		input := call.Argument(0).ToObject(r.vm)
		r.nativeScanStruct(r.vm, input)
		return goja.Undefined()
	})

	return nil
}

func (r *Runtime) nativePrepareRefs(source map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for sourceKey, value := range source {
		keyParts := strings.Split(sourceKey, ".")
		level := result

		for i, part := range keyParts {

			// If we are at the last part, assign the value
			if i == len(keyParts)-1 {
				level[part] = value
			} else {

				v, ok := level[part]

				// If the part does not exist, create a new map
				//if _, ok := level[part]; !ok {
				if !ok {
					v = make(map[string]interface{})
					level[part] = v
				}

				// Move deeper into the nested map
				nextLevel, _ := v.(map[string]interface{})
				level = nextLevel
			}
		}
	}

	return result
}

func (r *Runtime) nativeScanStruct(vm *goja.Runtime, obj *goja.Object) {

	// Get all keys of the object
	keys := obj.Keys()
	for _, key := range keys {
		value := obj.Get(key)

		if goja.IsUndefined(value) || goja.IsNaN(value) || goja.IsInfinity(value) {
			// Delete the key if the value is undefined
			obj.Delete(key)
		} else if goja.IsNull(value) {
			// Continue if the value is null
			continue
		} else if value.ExportType().Kind() == reflect.Map {
			// If the value is an object, recursively call scanStruct
			r.nativeScanStruct(vm, value.ToObject(vm))
		} else if value.ExportType().Kind() == reflect.Slice {
			// If the value is an array, iterate over the elements and recursively call scanStruct
			arrayObj := value.ToObject(vm)
			arrKey := arrayObj.Keys()
			for _, k := range arrKey {
				elem := arrayObj.Get(k)

				if goja.IsUndefined(elem) || goja.IsNaN(elem) || goja.IsInfinity(elem) || goja.IsNull(elem) {
					// Set undefined elements to null or delete them as per requirement
					arrayObj.Set(k, nil)
				} else {
					r.nativeScanStruct(vm, elem.ToObject(vm))
				}
			}
		}
	}
}
