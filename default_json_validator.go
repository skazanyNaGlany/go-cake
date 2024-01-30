package go_cake

import (
	"github.com/santhosh-tekuri/jsonschema/v5"
)

type DefaultJSONValidator struct {
	SchemaFilename string
	Schema         *jsonschema.Schema
}

func NewDefaultJSONValidator(schemaFilename string, schema string) (JSONValidator, error) {
	var err error

	validator := DefaultJSONValidator{SchemaFilename: schemaFilename}
	validator.Schema, err = jsonschema.CompileString(schemaFilename, schema)

	return &validator, err
}

func (djv *DefaultJSONValidator) Validate(item map[string]any) error {
	return djv.Schema.Validate(item)
}
