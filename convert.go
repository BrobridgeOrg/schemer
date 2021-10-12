package schemer

import (
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/BrobridgeOrg/schemer/types"
)

func getStandardValue(data interface{}) interface{} {

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

func getValue(def *Definition, data interface{}) interface{} {

	v := getStandardValue(data)

	// According to definition to convert value to what we want
	switch def.Type {
	case TYPE_INT64:
		return getIntegerValue(v)
	case TYPE_UINT64:
		return getUnsignedIntegerValue(v)
	case TYPE_FLOAT64:
		return getFloatValue(v)
	case TYPE_BOOLEAN:
		return getBoolValue(v)
	case TYPE_STRING:
		return getStringValue(v)
	case TYPE_TIME:
		return def.Info.(*types.Time).GetValue(v)
	case TYPE_BINARY:
		return getBinaryValue(v)
	}

	return v
}

func getIntegerValue(data interface{}) int64 {

	switch d := data.(type) {
	case int64:
		return d
	case uint64:
		return int64(d)
	case string:
		result, err := strconv.ParseInt(d, 10, 64)
		if err != nil {
			return 0
		}

		return result
	case bool:
		if d {
			return int64(1)
		} else {
			return int64(0)
		}
	case float64:
		return int64(d)
	case time.Time:
		return d.Unix()
	}

	return 0
}

func getUnsignedIntegerValue(data interface{}) uint64 {

	switch d := data.(type) {
	case int64:
		if d > 0 {
			return uint64(d)
		}

		return 0
	case uint64:
		return d
	case string:
		result, err := strconv.ParseUint(d, 10, 64)
		if err != nil {
			return 0
		}

		return result
	case bool:
		if d {
			return uint64(1)
		} else {
			return uint64(0)
		}
	case float64:
		return uint64(d)
	case time.Time:
		return uint64(d.Unix())
	}

	return 0
}

func getFloatValue(data interface{}) float64 {

	switch d := data.(type) {
	case int64:
		return float64(d)
	case uint64:
		return float64(d)
	case string:
		result, err := strconv.ParseFloat(d, 64)
		if err != nil {
			return 0
		}

		return result
	case bool:
		if d {
			return float64(1)
		} else {
			return float64(0)
		}
	case float64:
		return d
	case time.Time:
		return float64(d.Unix())
	}

	return 0
}

func getBoolValue(data interface{}) bool {

	switch d := data.(type) {
	case int64:
		if d > 0 {
			return true
		} else {
			return false
		}
	case uint64:
		if d > 0 {
			return true
		} else {
			return false
		}
	case string:
		result, err := strconv.ParseBool(d)
		if err != nil {
			return false
		}

		return result
	case bool:
		return d
	case float64:
		if d > 0 {
			return true
		} else {
			return false
		}
	case time.Time:
		return true
	}

	return false
}

func getStringValue(data interface{}) string {

	switch d := data.(type) {
	case int64:
		return fmt.Sprintf("%d", d)
	case uint64:
		return fmt.Sprintf("%d", d)
	case string:
		return d
	case bool:
		return fmt.Sprintf("%t", d)
	case float64:
		return strconv.FormatFloat(d, 'f', -1, 64)
	case time.Time:
		return d.UTC().Format(time.RFC3339Nano)
	default:
		return fmt.Sprintf("%v", d)
	}
}

func getBinaryValue(data interface{}) []byte {

	switch d := data.(type) {
	case []byte:
		return d
	case string:
		return []byte(d)
	default:

		arr, ok := data.([]interface{})
		if !ok {
			return []byte("")
		}

		val := make([]byte, len(arr))
		for i, v := range arr {
			val[i] = byte(getUnsignedIntegerValue(v))
		}

		return val
	}
}

func convert(sourceDef *Definition, destDef *Definition, data interface{}) interface{} {

	srcData := getValue(sourceDef, data)

	return getValue(destDef, srcData)
}
