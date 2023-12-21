package schemer

import "reflect"

type Record struct {
	schema *Schema
	raw    map[string]interface{}
}

func NewRecord(schema *Schema, raw map[string]interface{}) *Record {
	return &Record{
		schema: schema,
		raw:    raw,
	}
}

func (r *Record) GetValue(valuePath string) *Value {

	parts := r.schema.parsePath(valuePath)
	def := r.schema.getDefinition(parts)

	// Create a new value from raw data
	value := NewValue(def)
	value.Data = getValue(def, r.getValue(parts))

	return value
}

func (r *Record) getValue(parts []string) interface{} {

	var obj interface{} = r.raw
	var val interface{} = nil

	for _, p := range parts {

		if obj == nil {
			return nil
		}

		key, index := parsePathEntry(p)

		if v, ok := obj.(map[string]interface{}); ok {
			obj = v[key]
			val = obj

			if reflect.TypeOf(obj).Kind() == reflect.Slice && index != -1 {
				obj = obj.([]interface{})[index]
				val = obj
			}

			continue
		}

		val = obj
	}

	return val
}
