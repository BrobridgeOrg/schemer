package schemer

import (
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Schemer struct {
	schema *Schema
}

func NewSchemer() *Schemer {
	return &Schemer{}
}
