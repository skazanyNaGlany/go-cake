package go_cake

import (
	"fmt"
	"log"
	"runtime/debug"
	"strings"

	"github.com/skazanyNaGlany/go-cake/utils"
)

type BaseError struct {
	Message string
}

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

func (be BaseError) Error() string {
	return be.Message
}

func (be *BaseError) FormatStatusMessage(message string, baseError, internalError error) string {
	baseErrorType := utils.StructUtilsInstance.GetCleanType(baseError)
	baseErrorType = strings.Replace(baseErrorType, "Error", "", 1)

	formatted := fmt.Sprintf("%s: %s", baseErrorType, message)
	formatted = strings.TrimSpace(formatted)

	if internalError != nil {
		internalErrorMessage := strings.TrimSpace(internalError.Error())

		if internalErrorMessage != "" {
			if strings.LastIndex(formatted, ":") != len(formatted)-1 {
				formatted += ":"
			}

			formatted += fmt.Sprintf(" %s", internalErrorMessage)
		}
	}

	formatted = strings.TrimSpace(formatted)

	return strings.TrimSuffix(formatted, ":")
}

func (be *BaseError) logError(childErr error, object any) {
	message := fmt.Sprintf("HTTPError: %T\n", childErr)
	message += fmt.Sprintf("Stacktrace: %v\n", strings.TrimSpace(string(debug.Stack())))

	if object != nil {
		message += fmt.Sprintf("Object: %T (%v)\n", object, object)
	}

	log.Print(message)
}

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
	model GoKateModel,
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
	model GoKateModel,
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
	model GoKateModel,
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
