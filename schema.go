package schemer

import (
	"encoding/json"
	"strings"
	"sync"
)

type Schema struct {
	Fields map[string]*Definition
	Mutex  sync.RWMutex
}

func NewSchema() *Schema {
	return &Schema{
		Fields: make(map[string]*Definition),
	}
}

func (s *Schema) parsePath(fullPath string) []string {

	quoted := false
	elements := strings.FieldsFunc(fullPath, func(r rune) bool {

		if r == '"' {
			quoted = !quoted

			// Ignore
			return true
		}

		return !quoted && r == '.'
	})

	parts := make([]string, len(elements))
	for i, element := range elements {
		parts[i] = element
	}

	return parts
}

func (s *Schema) GetDefinition(valuePath string) *Definition {
	parts := s.parsePath(valuePath)
	return s.getDefinition(parts)
}

func (s *Schema) getDefinition(parts []string) *Definition {

	s.Mutex.RLock()
	defer s.Mutex.RUnlock()

	var def *Definition
	fields := s.Fields
	for _, path := range parts {
		d, ok := fields[path]
		if !ok {
			return nil
		}

		def = d

		if d.Type == TYPE_MAP {
			fields = def.Definition.Fields
		}
	}

	return def
}

func (s *Schema) normalize(schema *Schema, data map[string]interface{}) map[string]interface{} {

	result := make(map[string]interface{})

	for fieldName, def := range schema.Fields {

		val, ok := data[fieldName]
		if !ok {
			continue
		}

		if def.Type == TYPE_MAP {
			result[fieldName] = s.normalize(def.Definition, val.(map[string]interface{}))
			continue
		}

		result[fieldName] = getValue(def, val)
	}

	return result
}

func (s *Schema) Normalize(data map[string]interface{}) map[string]interface{} {
	return s.normalize(s, data)
}

func (s *Schema) Scan(data map[string]interface{}) *Record {
	return NewRecord(s, s.normalize(s, data))
}

func UnmarshalJSON(source []byte, s *Schema) error {

	// Parsing original JSON string
	var raw map[string]interface{}
	err := json.Unmarshal([]byte(source), &raw)
	if err != nil {
		return err
	}

	// Create a new schema
	err = Unmarshal(raw, s)
	if err != nil {
		return err
	}

	return nil
}

func Unmarshal(data map[string]interface{}, s *Schema) error {

	for key, value := range data {

		// Parse definition from unknown interface object
		var def Definition
		err := UnmarshalDefinition(value, &def)
		if err != nil {
			return err
		}

		s.Fields[key] = &def
	}

	return nil
}
