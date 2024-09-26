package schemer

import (
	"errors"

	"github.com/BrobridgeOrg/schemer/types"
)

var (
	ErrInvalidTypeDefinition    = errors.New("Invalid type definition")
	ErrInvalidFieldsDefinition  = errors.New("Invalid fields definition")
	ErrInvalidNotNullDefinition = errors.New("Invalid notNull definition")
	ErrInvalidArraySubtype      = errors.New("Array type requires subtype")
)

type RawDefinition struct {
	Type    ValueType
	Subtype *RawDefinition
	Fields  map[string]*RawDefinition
	NotNull bool
	Props   map[string]interface{}
}

type Definition struct {
	Schema  *Schema
	Type    ValueType
	Subtype *Definition
	Info    interface{}
	NotNull bool
}

func NewRawDefinition() *RawDefinition {
	return &RawDefinition{
		Type:    TYPE_ANY,
		NotNull: false,
		Props:   make(map[string]interface{}),
	}
}

func NewDefinition(t ValueType) *Definition {
	return &Definition{
		Schema:  nil,
		Type:    t,
		Subtype: nil,
	}
}

func UnmarshalDefinition(data interface{}, d *Definition) error {

	raw, err := extractRawDefinition(data)
	if err != nil {
		return err
	}

	def, err := createDefinitionFromRawDefinition(raw)
	if err != nil {
		return err
	}

	// Fill the definition
	d.Type = def.Type
	d.Subtype = def.Subtype
	d.Schema = def.Schema
	d.Info = def.Info
	d.NotNull = def.NotNull

	return nil
}

func createDefinitionFromRawDefinition(raw *RawDefinition) (*Definition, error) {

	def := NewDefinition(raw.Type)
	def.NotNull = raw.NotNull

	switch def.Type {
	case TYPE_MAP:
		s, err := createSchemaFromRawFields(raw.Fields)
		if err != nil {
			return nil, err
		}

		def.Schema = s

	case TYPE_ARRAY:

		if raw.Subtype == nil {
			return nil, ErrInvalidArraySubtype
		}

		subDef, err := createDefinitionFromRawDefinition(raw.Subtype)
		if err != nil {
			return nil, err
		}

		def.Subtype = subDef

	case TYPE_TIME:
		t := types.NewTime()
		t.Parse(raw.Props)
		def.Info = t
	}

	return def, nil
}

func createSchemaFromRawFields(rawMap map[string]*RawDefinition) (*Schema, error) {

	s := NewSchema()

	for key, value := range rawMap {
		raw, err := createDefinitionFromRawDefinition(value)
		if err != nil {
			return nil, err
		}
		s.Fields[key] = raw
	}

	return s, nil
}

func extractSubtypeRawDefinition(parent map[string]interface{}, v interface{}) (*RawDefinition, error) {

	switch d := v.(type) {
	case string:
		return extractRawDefinition(map[string]interface{}{
			"type":   d,
			"fields": parent["fields"],
		})
	case map[string]interface{}:
		return extractRawDefinition(d)
	}

	return nil, ErrInvalidArraySubtype
}

func extractFieldsRawDefinition(v interface{}) (map[string]*RawDefinition, error) {

	switch d := v.(type) {
	case map[string]interface{}:

		fields := make(map[string]*RawDefinition)

		for key, value := range d {

			raw, err := extractRawDefinition(value)
			if err != nil {
				return nil, err
			}

			fields[key] = raw
		}

		return fields, nil
	}

	return nil, ErrInvalidFieldsDefinition
}

func extractRawDefinition(data interface{}) (*RawDefinition, error) {

	raw := NewRawDefinition()

	switch v := data.(type) {
	case map[string]interface{}:

		// Handle type
		t, ok := v["type"]
		if !ok {
			return nil, ErrInvalidTypeDefinition
		}

		switch d := t.(type) {
		case string:

			vt, ok := ValueTypes[d]
			if !ok {
				return nil, ErrInvalidType
			}

			raw.Type = vt

		default:
			return nil, ErrInvalidTypeDefinition
		}

		// Handle subtype for array
		if raw.Type == TYPE_ARRAY {

			st, ok := v["subtype"]
			if !ok {
				return nil, ErrInvalidArraySubtype
			}

			// Decode subtype
			subDef, err := extractSubtypeRawDefinition(v, st)
			if err != nil {
				return nil, err
			}

			raw.Subtype = subDef
		}

		// Handle fields
		f, ok := v["fields"]
		if raw.Type == TYPE_MAP && (!ok || f == nil) {

			// Fields is required for map type
			return nil, ErrInvalidFieldsDefinition
		}

		if ok && f != nil {
			// Attempt to extract fields
			fields, err := extractFieldsRawDefinition(f)
			if err != nil {
				return nil, err
			}

			if raw.Type == TYPE_MAP {
				raw.Fields = fields
			}

			// Compatible with old version, which uses subtype.fields instead of fields
			if raw.Type == TYPE_ARRAY && len(raw.Fields) > 0 &&
				(raw.Subtype != nil && raw.Subtype.Type == TYPE_MAP && raw.Subtype.Fields == nil) {
				raw.Subtype.Fields = fields
			}
		}

		if val, ok := v["notNull"]; ok {
			switch val.(type) {
			case bool:
				raw.NotNull = val.(bool)
			default:
				return nil, ErrInvalidNotNullDefinition
			}
		}

		// More properties
		for key, value := range v {

			// Skip reserved keys
			switch key {
			case "type", "fields", "notNull":
				continue
			}

			raw.Props[key] = value
		}

	}

	return raw, nil
}
