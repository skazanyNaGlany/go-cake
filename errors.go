package go_cake

import (
	"fmt"
)

type RestHandlerPatternExistsError struct{ BaseError }
type NoResourceDatabaseDriverSetError struct{ BaseError }
type NoResourceDbPathSetError struct{ BaseError }
type NoResourcePatternSetError struct{ BaseError }
type NoResourceResourceNameSetError struct{ BaseError }
type NoResourceDbModelSetError struct{ BaseError }
type NoResourceJSONSchemaConfigSetError struct{ BaseError }
type NoResourceJSONSchemaConfigIDFieldSetError struct{ BaseError }
type IDFieldModelNotFoundError struct{ BaseError }
type ETagFieldModelNotFoundError struct{ BaseError }
type UnableToTestModelError struct{ BaseError }
type UnableToInitDatabaseDriverError struct{ BaseError }
type NoSchemaConfigError struct{ BaseError }
type NoSchemaConfigIDError struct{ BaseError }
type SchemaConfigUnknownFieldError struct{ BaseError }

func NewNoResourceDatabaseDriverSetError(resource *Resource, internalError error) error {
	e := NoResourceDatabaseDriverSetError{}

	e.Message = e.FormatStatusMessage(
		fmt.Sprintf("no DatabaseDriver set for %T resource", resource),
		e,
		internalError)

	e.logError(e, nil)

	return e
}

func NewNoResourceDbPathSetError(resource *Resource, internalError error) error {
	e := NoResourceDbPathSetError{}

	e.Message = e.FormatStatusMessage(
		fmt.Sprintf("no DbPath set for %T resource", resource),
		e,
		internalError)

	e.logError(e, nil)

	return e
}

func NewNoResourcePatternSetError(resource *Resource, internalError error) error {
	e := NoResourcePatternSetError{}

	e.Message = e.FormatStatusMessage(
		fmt.Sprintf("no Pattern set for %T resource", resource),
		e,
		internalError)

	e.logError(e, nil)

	return e
}

func NewNoResourceResourceNameSetError(resource *Resource, internalError error) error {
	e := NoResourceResourceNameSetError{}

	e.Message = e.FormatStatusMessage(
		fmt.Sprintf("no ResourceName set for %T resource", resource),
		e,
		internalError)

	e.logError(e, nil)

	return e
}

func NewNoResourceDbModelSetError(resource *Resource, internalError error) error {
	e := NoResourceDbModelSetError{}

	e.Message = e.FormatStatusMessage(
		fmt.Sprintf("no DbModel set for %T resource", resource),
		e,
		internalError)

	e.logError(e, nil)

	return e
}

func NewNoResourceJSONSchemaConfigSetError(resource *Resource, internalError error) error {
	e := NoResourceJSONSchemaConfigSetError{}

	e.Message = e.FormatStatusMessage(
		fmt.Sprintf("no JSONSchemaConfig set for %T resource", resource),
		e,
		internalError)

	e.logError(e, nil)

	return e
}

func NewNoResourceJSONSchemaConfigIDFieldSetError(resource *Resource, internalError error) error {
	e := NoResourceJSONSchemaConfigIDFieldSetError{}

	e.Message = e.FormatStatusMessage(
		fmt.Sprintf("no JSONSchemaConfig.IDField set for %T resource", resource),
		e,
		internalError)

	e.logError(e, nil)

	return e
}

func NewIDFieldModelNotFoundError(
	resource *Resource,
	model GoCakeModel,
	jsonIDField string,
	internalError error) error {
	e := IDFieldModelNotFoundError{}

	e.Message = e.FormatStatusMessage(
		fmt.Sprintf(
			"JSON ID field %v not found in %T model for %T resource",
			jsonIDField,
			model,
			resource),
		e,
		internalError)

	e.logError(e, nil)

	return e
}

func NewETagFieldModelNotFoundError(
	resource *Resource,
	model GoCakeModel,
	jsonETagField string,
	internalError error) error {
	e := ETagFieldModelNotFoundError{}

	e.Message = e.FormatStatusMessage(
		fmt.Sprintf(
			"JSON ETag field %v not found in %T model for %T resource",
			jsonETagField,
			model,
			resource),
		e,
		internalError)

	e.logError(e, nil)

	return e
}

func NewUnableToTestModelError(
	resource *Resource,
	driver DatabaseDriver,
	model GoCakeModel,
	internalError error) error {
	e := UnableToTestModelError{}

	e.Message = e.FormatStatusMessage(
		fmt.Sprintf(
			"Unable to test model %T against %T driver for %T resource",
			model,
			driver,
			resource),
		e,
		internalError)

	e.logError(e, nil)

	return e
}

// TODO add messages to each error
