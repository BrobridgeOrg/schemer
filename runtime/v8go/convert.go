package v8go_runtime

import (
	"fmt"
	"log"
	"reflect"

	msgpack "github.com/vmihailenco/msgpack/v5"

	"rogchap.com/v8go"
)

func (r *Runtime) normalizeValue(data interface{}) interface{} {

	v := reflect.ValueOf(data)

	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint()
	case reflect.Float32, reflect.Float64:
		return v.Float()
	}

	return data
}

func (r *Runtime) convertToV8Object(data interface{}) *v8go.Value {
	/*
		jsonString, _ := json.Marshal(data)

		r.ctx.Global().Set("input", string(jsonString))
	*/

	if data == nil {
		return v8go.Null(r.iso)
	}

	encodedData, err := msgpack.Marshal(data)
	if err != nil {
		return nil
	}

	script := fmt.Sprintf("new Uint8Array(%d)", len(encodedData))

	val, err := r.ctx.RunScript(script, "arraybuffer.js")
	if err != nil {
		panic(err)
	}

	obj := val.Object()

	for index, value := range encodedData {

		v8Value, err := v8go.NewValue(r.iso, int32(value))
		if err != nil {
			log.Fatalf("Error creating v8 value for byte %d: %v", value, err)
		}

		obj.SetIdx(uint32(index), v8Value)
	}

	r.ctx.Global().Set("inputData", obj)
	returnedValue, err := r.ctx.RunScript("msgpack.decode(inputData)", "result.js")
	if err != nil {
		return nil
	}

	return returnedValue
}

func (r *Runtime) convertV8ObjectToStruct(value *v8go.Value) (interface{}, error) {

	// Put object into global object
	r.ctx.Global().Set("outputData", value.Object())
	//r.ctx.RunScript("console.log('stringify', JSON.stringify(result))", "json.js")
	returnedData, err := r.ctx.RunScript("msgpack.encode(outputData)", "result.js")
	if err != nil {
		return nil, err
	}

	// Convert maskpack to ][byte for golang
	source := returnedData.Object()

	length, err := source.Get("length")
	if err != nil {
		return nil, err
	}

	size := int(length.Int32())

	data := make([]byte, size)
	for i := 0; i < size; i++ {

		b, err := source.GetIdx(uint32(i))
		if err != nil {
			return nil, err
		}

		data[i] = byte(b.Int32())

		//fmt.Println(i, data)
	}

	var result interface{}
	err = msgpack.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	switch d := result.(type) {
	case map[string]interface{}:
		r.convertAllBufferToBytes(d)
	case []interface{}:
		for _, ele := range d {
			switch dd := ele.(type) {
			case map[string]interface{}:
				r.convertAllBufferToBytes(dd)
			}
		}
	}

	return result, nil
}

func (r *Runtime) convertAllBufferToBytes(data map[string]interface{}) interface{} {

	for key, value := range data {
		switch v := value.(type) {
		case map[string]interface{}:
			data[key] = r.convertAllBufferToBytes(v)
		case []interface{}:
			for i, ele := range v {
				switch vv := ele.(type) {
				case map[string]interface{}:
					v[i] = r.convertAllBufferToBytes(vv)
				}
			}
		case *Uint8Array:
			data[key] = v.Buffer.Bytes()
		}
	}

	return data
}
