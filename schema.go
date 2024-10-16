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
	for _, entry := range parts {

		if def != nil {

			switch def.Type {
			case TYPE_ARRAY:

				if def.Subtype == nil {
					return nil
				}

				def = def.Subtype

				if def.Type == TYPE_MAP {
					fields = def.Schema.Fields
				} else if def.Type == TYPE_ARRAY {
					continue
				} else {
					return def
				}
			case TYPE_MAP:
				fields = def.Schema.Fields
			}
		}

		// Parse key and index
		key, _ := parsePathEntry(entry)

		// Check if we have a definition for this key
		d, ok := fields[key]
		if !ok {
			// No definition found
			return nil
		}

		def = d
	}

	return def
}

func (s *Schema) normalize(schema *Schema, data map[string]interface{}) map[string]interface{} {

	result := make(map[string]interface{})

	for fieldName, def := range schema.Fields {

		// Skip internal fields
		if strings.HasPrefix(fieldName, "$") {
			continue
		}

		// Check if field name contains a path. If so, we need to parse it to check if the key exists.
		if strings.Contains(fieldName, ".") {
			pathDef := s.GetDefinition(fieldName)
			if pathDef == nil {
				// Skip this field if the path does not exist in the schema
				continue
			}
		}

		val, ok := data[fieldName]
		if !ok {
			continue
		}

		if def.Type == TYPE_MAP && val != nil {
			result[fieldName] = s.normalize(def.Schema, val.(map[string]interface{}))
			continue
		}

		v, _ := getValue(def, val)

		result[fieldName] = v
	}

	// set value by path
	for key, val := range data {

		// Keep internal fields
		if strings.HasPrefix(key, "$") {
			result[key] = val
			continue
		}

		// Check if field name contains a path.
		if !strings.Contains(key, ".") {
			continue
		}

		// We need to parse it to check if the key exists.
		def := s.GetDefinition(key)
		if def == nil {
			// Skip this field if the path does not exist in the schema
			continue
		}

		v, err := getValue(def, val)
		if err != nil {
			continue
		}

		result[key] = v
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
