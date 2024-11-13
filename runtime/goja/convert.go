package goja_runtime

import (
	"reflect"
	"time"

	"github.com/BrobridgeOrg/schemer"
	"github.com/BrobridgeOrg/schemer/types"
	"github.com/dop251/goja"
)

func (r *Runtime) handleMapValue(returnedValue map[string]interface{}) error {

	for key, value := range returnedValue {
		v := reflect.ValueOf(value)
		switch v.Kind() {
		/*
			case reflect.Slice:
				err := t.handleArrayValue(value.([]interface{}))
				if err != nil {
					return err
				}
		*/
		case reflect.Map:
			err := r.handleMapValue(value.(map[string]interface{}))
			if err != nil {
				return err
			}
		default:
			// Convert Data object to time.Time
			switch d := value.(type) {
			case *goja.Object:
				if value.(*goja.Object).ClassName() == "Date" {
					returnedValue[key] = d.Export()
				}
			}
		}
	}

	return nil
}

func (r *Runtime) normalize(schema *schemer.Schema, data map[string]interface{}) {

	for fieldName, def := range schema.Fields {

		val, ok := data[fieldName]
		if !ok {
			continue
		}

		if def.Type == schemer.TYPE_MAP {
			r.normalize(def.Schema, val.(map[string]interface{}))
			continue
		}

		if def.Type == schemer.TYPE_TIME {

			// Skip null
			if val == nil {
				continue
			}

			if def.Info.(*types.Time).Precision != types.TIME_PRECISION_MICROSECOND {

				// New Date object
				dateValue := val.(time.Time).UnixMicro() / 1e3
				v, _ := r.vm.New(r.vm.Get("Date").ToObject(r.vm), r.vm.ToValue(dateValue))
				data[fieldName] = v
			}
			continue
		}
	}
}
