package go_cake

import (
	"strings"

	"github.com/thoas/go-funk"
)

type JSONSchemaConfig struct {
	IDField                string
	ETagField              string
	FilterableFields       []string
	ProjectableFields      []string
	SortableFields         []string
	InsertableFields       []string
	UpdatableFields        []string
	HiddenFields           []string
	ErasedFields           []string
	RequiredOnInsertFields []string
	RequiredOnUpdateFields []string
	OptimizeOnInsertFields []string
	OptimizeOnUpdateFields []string
	Validator              JSONValidator
}

func (jsc *JSONSchemaConfig) GetAllFields() []string {
	allFields := make([]string, 0)
	filtered := make([]string, 0)

	allFields = append(allFields, jsc.IDField)
	allFields = append(allFields, jsc.ETagField)
	allFields = append(allFields, jsc.FilterableFields...)
	allFields = append(allFields, jsc.ProjectableFields...)
	allFields = append(allFields, jsc.SortableFields...)
	allFields = append(allFields, jsc.InsertableFields...)
	allFields = append(allFields, jsc.UpdatableFields...)
	allFields = append(allFields, jsc.HiddenFields...)
	allFields = append(allFields, jsc.ErasedFields...)
	allFields = append(allFields, jsc.RequiredOnInsertFields...)
	allFields = append(allFields, jsc.RequiredOnUpdateFields...)
	allFields = append(allFields, jsc.OptimizeOnInsertFields...)
	allFields = append(allFields, jsc.OptimizeOnUpdateFields...)

	allFields = funk.UniqString(allFields)

	for _, iField := range allFields {
		iField = strings.TrimSpace(iField)

		if iField == "" {
			continue
		}

		filtered = append(filtered, iField)
	}

	return filtered
}
