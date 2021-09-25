package schemer

type ValueType int32

const (
	TYPE_BOOLEAN ValueType = 0
	TYPE_BINARY  ValueType = 1
	TYPE_STRING  ValueType = 2
	TYPE_UINT64  ValueType = 3
	TYPE_INT64   ValueType = 4
	TYPE_FLOAT64 ValueType = 5
	TYPE_ARRAY   ValueType = 6
	TYPE_MAP     ValueType = 7
	TYPE_TIME    ValueType = 8
	TYPE_NULL    ValueType = 9
	TYPE_ANY     ValueType = 10
)

var ValueTypes = map[string]ValueType{
	"string": TYPE_STRING,
	"binary": TYPE_BINARY,
	"int":    TYPE_INT64,
	"uint":   TYPE_UINT64,
	"float":  TYPE_FLOAT64,
	"bool":   TYPE_BOOLEAN,
	"time":   TYPE_TIME,
	"array":  TYPE_ARRAY,
	"map":    TYPE_MAP,
	"any":    TYPE_ANY,
}

type Value struct {
	Definition *Definition
	Data       interface{}
}

func NewValue(def *Definition) *Value {
	return &Value{
		Definition: def,
	}
}
