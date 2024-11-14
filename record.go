package schemer

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
	if def == nil {
		return nil
	}

	// Create a new value from raw data
	value := NewValue(def)

	// get value with defintion type from raw data
	v, err := getValue(def, r.getValue(parts))
	if err != nil {
		return nil
	}

	value.Data = v

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

			if d, ok := obj.([]interface{}); ok && index != -1 {
				obj = d[index]
				val = d
			}

			/*
				if reflect.TypeOf(obj).Kind() == reflect.Slice && index != -1 {
					obj = obj.([]interface{})[index]
					val = obj
				}
			*/
			continue
		}

		val = obj
	}

	return val
}
