package go_cake

import (
	"fmt"
	"net/http"
	"strings"
)

type HTTPError interface {
	Error() string
	GetStatusCode() int
	GetStatusMessage() string
}

type ModifiersNotAllowedHTTPError struct{ BaseHTTPError }
type MethodNotAllowedHTTPError struct{ BaseHTTPError }
type UnauthorizedHTTPError struct{ BaseHTTPError }
type URLNotFoundHTTPError struct{ BaseHTTPError }
type UnableToParseRequestHTTPError struct{ BaseHTTPError }
type TooManyInputItemsHTTPError struct{ BaseHTTPError }
type URLTooBigHTTPError struct{ BaseHTTPError }
type PayloadTooBigHTTPError struct{ BaseHTTPError }
type PerPageTooLargeHTTPError struct{ BaseHTTPError }
type OKHTTPError struct{ BaseHTTPError }
type FieldNotExistsHTTPError struct{ BaseHTTPError }
type FieldNotFilterableHTTPError struct{ BaseHTTPError }
type FieldNotSortableHTTPError struct{ BaseHTTPError }
type FieldNotProjectableHTTPError struct{ BaseHTTPError }
type FieldRequiredHTTPError struct{ BaseHTTPError }
type FieldNotInsertableHTTPError struct{ BaseHTTPError }
type FieldNotUpdatableHTTPError struct{ BaseHTTPError }
type InvalidAcceptRequestHeaderHTTPError struct{ BaseHTTPError }
type InvalidContentTypeRequestHeaderHTTPError struct{ BaseHTTPError }
type ClientObjectMalformedHTTPError struct{ BaseHTTPError }
type ServerObjectMalformedHTTPError struct{ BaseHTTPError }
type PayloadInvalidHTTPError struct{ BaseHTTPError }
type CannotDecodePayloadHTTPError struct{ BaseHTTPError }
type LowLevelDriverHTTPError struct{ BaseHTTPError }
type InternalServerErrorHTTPError struct{ BaseHTTPError }
type MalformedWhereHTTPError struct{ BaseHTTPError }
type MalformedSortHTTPError struct{ BaseHTTPError }
type MalformedProjectionHTTPError struct{ BaseHTTPError }
type ObjectNotFoundHTTPError struct{ BaseHTTPError }
type ObjectNotAffectedHTTPError struct{ BaseHTTPError }
type TooManyAffectedObjectsHTTPError struct{ BaseHTTPError }
type UnsupportedVersionHTTPError struct{ BaseHTTPError }

func NewMethodNotAllowedHTTPError(internalError error) HTTPError {
	e := MethodNotAllowedHTTPError{}

	e.StatusCode = http.StatusMethodNotAllowed
	e.StatusMessage = e.FormatStatusMessage(
		http.StatusText(http.StatusMethodNotAllowed),
		e,
		internalError)

	return e
}

func NewUnauthorizedHTTPError(internalError error) HTTPError {
	e := UnauthorizedHTTPError{}

	e.StatusCode = http.StatusUnauthorized
	e.StatusMessage = e.FormatStatusMessage("", e, internalError)

	return e
}

func NewURLNotFoundHTTPError(internalError error) HTTPError {
	e := URLNotFoundHTTPError{}

	e.StatusCode = http.StatusNotFound
	e.StatusMessage = e.FormatStatusMessage("URL not found", e, internalError)

	return e
}

func NewUnableToParseRequestHTTPError(internalError error) HTTPError {
	e := UnableToParseRequestHTTPError{}

	e.StatusCode = http.StatusInternalServerError
	e.StatusMessage = e.FormatStatusMessage("Unable to parse request", e, internalError)

	e.logError(e, nil)

	return e
}

func NewTooManyInputItemsHTTPError(maxItems int64, gotItems int, internalError error) HTTPError {
	e := TooManyInputItemsHTTPError{}

	message := fmt.Sprintf("Number of the items exceeded maximum limit of %v (got %v items)", maxItems, gotItems)

	e.StatusCode = http.StatusBadRequest
	e.StatusMessage = e.FormatStatusMessage(message, e, internalError)

	return e
}

func NewURLTooBigHTTPError(maxLenght uint64, internalError error) HTTPError {
	e := URLTooBigHTTPError{}

	message := fmt.Sprintf("URL length exceeded maximum limit of %d bytes", maxLenght)

	e.StatusCode = http.StatusRequestURITooLong
	e.StatusMessage = e.FormatStatusMessage(message, e, internalError)

	return e
}

func NewPayloadTooBigHTTPError(maxSize int64, internalError error) HTTPError {
	e := PayloadTooBigHTTPError{}

	message := fmt.Sprintf("Payload size exceeded maximum limit of %d bytes", maxSize)

	e.StatusCode = http.StatusRequestEntityTooLarge
	e.StatusMessage = e.FormatStatusMessage(message, e, internalError)

	return e
}

func NewPerPageTooLargeHTTPError(max int64, internalError error) HTTPError {
	e := PerPageTooLargeHTTPError{}

	message := fmt.Sprintf("Maximum allowed per_page value is %d", max)

	e.StatusCode = http.StatusRequestEntityTooLarge
	e.StatusMessage = e.FormatStatusMessage(message, e, internalError)

	return e
}

func NewOKHTTPError(internalError error) HTTPError {
	e := OKHTTPError{}

	e.StatusCode = http.StatusOK
	e.StatusMessage = e.FormatStatusMessage("", e, internalError)

	return e
}

func NewFieldNotExistsHTTPError(field string, internalError error) HTTPError {
	e := FieldNotExistsHTTPError{}

	message := fmt.Sprintf("Field '%v' does not exists", field)

	e.StatusCode = http.StatusBadRequest
	e.StatusMessage = e.FormatStatusMessage(message, e, internalError)

	return e
}

func NewFieldNotFilterableHTTPError(field string, internalError error) HTTPError {
	e := FieldNotFilterableHTTPError{}

	message := fmt.Sprintf("Field '%v' is not filterable", field)

	e.StatusCode = http.StatusBadRequest
	e.StatusMessage = e.FormatStatusMessage(message, e, internalError)

	return e
}

func NewFieldNotSortableHTTPError(field string, internalError error) HTTPError {
	e := FieldNotSortableHTTPError{}

	message := fmt.Sprintf("Field '%v' is not sortable", field)

	e.StatusCode = http.StatusBadRequest
	e.StatusMessage = e.FormatStatusMessage(message, e, internalError)

	return e
}

func NewFieldNotProjectableHTTPError(field string, internalError error) HTTPError {
	e := FieldNotProjectableHTTPError{}

	message := fmt.Sprintf("Field '%v' is not projectable", field)

	e.StatusCode = http.StatusBadRequest
	e.StatusMessage = e.FormatStatusMessage(message, e, internalError)

	return e
}

func NewFieldRequiredHTTPError(field string, internalError error) HTTPError {
	e := FieldRequiredHTTPError{}

	message := fmt.Sprintf("Field '%v' is required", field)

	e.StatusCode = http.StatusBadRequest
	e.StatusMessage = e.FormatStatusMessage(message, e, internalError)

	return e
}

func NewFieldNotInsertableHTTPError(field string, internalError error) HTTPError {
	e := FieldNotInsertableHTTPError{}

	message := fmt.Sprintf("Field '%v' is not insertable", field)

	e.StatusCode = http.StatusBadRequest
	e.StatusMessage = e.FormatStatusMessage(message, e, internalError)

	return e
}

func NewFieldNotUpdatableHTTPError(field string, internalError error) HTTPError {
	e := FieldNotUpdatableHTTPError{}

	message := fmt.Sprintf("Field '%v' is not updatable", field)

	e.StatusCode = http.StatusBadRequest
	e.StatusMessage = e.FormatStatusMessage(message, e, internalError)

	return e
}

func NewInvalidAcceptRequestHeaderHTTPError(allowed []string, internalError error) HTTPError {
	e := InvalidAcceptRequestHeaderHTTPError{}

	message := fmt.Sprintf("Invalid or missing Accept request header; allowed values are %v", strings.Join(allowed, ", "))

	e.StatusCode = http.StatusBadRequest
	e.StatusMessage = e.FormatStatusMessage(message, e, internalError)

	return e
}

func NewInvalidContentTypeRequestHeaderHTTPError(allowed string, internalError error) HTTPError {
	e := InvalidContentTypeRequestHeaderHTTPError{}

	message := fmt.Sprintf("Invalid or missing Content-Type request header; allowed values are %v", allowed)

	e.StatusCode = http.StatusUnsupportedMediaType
	e.StatusMessage = e.FormatStatusMessage(message, e, internalError)

	return e
}

func NewClientObjectMalformedHTTPError(internalError error) HTTPError {
	e := ClientObjectMalformedHTTPError{}

	message := "Client object is malformed"

	e.StatusCode = http.StatusBadRequest
	e.StatusMessage = e.FormatStatusMessage(message, e, internalError)

	return e
}

func NewServerObjectMalformedHTTPError(object any, internalError error) HTTPError {
	e := ServerObjectMalformedHTTPError{}

	message := "Server object is malformed"

	e.StatusCode = http.StatusInternalServerError
	e.StatusMessage = e.FormatStatusMessage(message, e, internalError)

	e.logError(e, object)

	return e
}

func NewPayloadInvalidHTTPError(internalError error) HTTPError {
	e := PayloadInvalidHTTPError{}

	e.StatusCode = http.StatusBadRequest
	e.StatusMessage = e.FormatStatusMessage("Passed payload is invalid, cannot be decoded or contains errors", e, internalError)

	return e
}

func NewCannotDecodePayloadHTTPError(internalError error) HTTPError {
	e := CannotDecodePayloadHTTPError{}

	e.StatusCode = http.StatusBadRequest
	e.StatusMessage = e.FormatStatusMessage("Passed payload cannot be decoded", e, internalError)

	return e
}

func NewLowLevelDriverHTTPError(internalError error) HTTPError {
	e := LowLevelDriverHTTPError{}

	e.StatusCode = http.StatusInternalServerError
	e.StatusMessage = e.FormatStatusMessage("", e, internalError)

	e.logError(e, nil)

	return e
}

func NewInternalServerErrorHTTPError(internalError error) HTTPError {
	e := InternalServerErrorHTTPError{}

	e.StatusCode = http.StatusInternalServerError
	e.StatusMessage = e.FormatStatusMessage("", e, internalError)

	e.logError(e, nil)

	return e
}

func NewMalformedWhereHTTPError(internalError error) HTTPError {
	e := MalformedWhereHTTPError{}

	e.StatusCode = http.StatusBadRequest
	e.StatusMessage = e.FormatStatusMessage("", e, internalError)

	return e
}

func NewMalformedSortHTTPError(internalError error) HTTPError {
	e := MalformedSortHTTPError{}

	e.StatusCode = http.StatusBadRequest
	e.StatusMessage = e.FormatStatusMessage("", e, internalError)

	return e
}

func NewMalformedProjectionHTTPError(internalError error) HTTPError {
	e := MalformedProjectionHTTPError{}

	e.StatusCode = http.StatusBadRequest
	e.StatusMessage = e.FormatStatusMessage("", e, internalError)

	return e
}

func NewObjectNotFoundHTTPError(internalError error) HTTPError {
	e := ObjectNotFoundHTTPError{}

	e.StatusCode = http.StatusNotFound
	e.StatusMessage = e.FormatStatusMessage("Object not found", e, internalError)

	return e
}

func NewObjectNotAffectedHTTPError(internalError error) HTTPError {
	e := ObjectNotAffectedHTTPError{}

	e.StatusCode = http.StatusInternalServerError
	e.StatusMessage = e.FormatStatusMessage("Object not affected (not inserted, not updated, not deleted)", e, internalError)

	return e
}

func NewTooManyAffectedObjectsHTTPError(internalError error) HTTPError {
	e := TooManyAffectedObjectsHTTPError{}

	e.StatusCode = http.StatusRequestEntityTooLarge
	e.StatusMessage = e.FormatStatusMessage("Affected more than one object with the same ID", e, internalError)

	return e
}

func NewUnsupportedVersionHTTPError(version string, internalError error) HTTPError {
	e := UnsupportedVersionHTTPError{}

	message := fmt.Sprintf("Passed API version '%s' is not supported", version)

	e.StatusCode = http.StatusNotAcceptable
	e.StatusMessage = e.FormatStatusMessage(message, e, internalError)

	return e
}

func NewModifiersNotAllowedHTTPError(internalError error) HTTPError {
	e := ModifiersNotAllowedHTTPError{}

	e.StatusCode = http.StatusBadRequest
	e.StatusMessage = e.FormatStatusMessage("Modifiers not allowed", e, internalError)

	return e
}
