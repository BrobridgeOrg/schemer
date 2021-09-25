package schemer

import (
	"fmt"

	"github.com/BrobridgeOrg/schemer/types"
	"github.com/mitchellh/mapstructure"
)

type RawDefinition struct {
	Type    string                 `mapstructure:"type"`
	Subtype string                 `mapstructure:"subtype"`
	Fields  map[string]interface{} `mapstructure:"fields"`
}

type Definition struct {
	Definition *Schema
	Type       ValueType
	Subtype    ValueType
	Info       interface{}
}

func NewDefinition(t ValueType) *Definition {
	return &Definition{
		Definition: nil,
		Type:       t,
	}
}

func UnmarshalDefinition(data interface{}, d *Definition) error {

	var raw RawDefinition
	err := mapstructure.Decode(data, &raw)
	if err != nil {
		return fmt.Errorf("Unknown defninition: %v", data)
	}

	d.Type = ValueTypes[raw.Type]
	switch d.Type {
	case TYPE_ARRAY:
		d.Subtype = ValueTypes[raw.Subtype]
	case TYPE_MAP:
		s := NewSchema()
		err = Unmarshal(raw.Fields, s)
		if err != nil {
			return err
		}

		d.Definition = s
	case TYPE_TIME:
		t := types.NewTime()
		t.Parse(data)
		d.Info = t
	}

	return nil
}
