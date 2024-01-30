package go_cake

import (
	"github.com/thoas/go-funk"
)

type UpdateRequestProcessor struct {
	BaseRequestProcessor
}

func NewUpdateRequestProcessor(request *Request, resource *Resource) *UpdateRequestProcessor {
	var updateRequestProcessor UpdateRequestProcessor

	updateRequestProcessor.request = request
	updateRequestProcessor.resource = resource
	updateRequestProcessor.subRequestProcessor = &updateRequestProcessor

	return &updateRequestProcessor
}

// TODO require _id and _etag in the payload
// when they are set in the resource
func (urp *UpdateRequestProcessor) ProcessRequest(response *ResponseJSON) ([]GoKateModel, HTTPError) {
	var httpErr HTTPError

	if !urp.resource.UpdateAllowed {
		return nil, NewMethodNotAllowedHTTPError(nil)
	}

	if httpErr = urp.checkRanges(); httpErr != nil {
		return nil, httpErr
	}

	if urp.request.HasWhere() || urp.request.HasSort() || urp.request.HasPage() {
		return nil, NewModifiersNotAllowedHTTPError(nil)
	}

	urp.optimizeFields(urp.request.DecodedJsonSlice)
	urp.preRequestJSONActions(urp.request.DecodedJsonSlice)

	converted, err := urp.decodedJsonSliceToDBModels(urp.request.DecodedJsonSlice)

	if err != nil {
		return converted, err
	}

	httpErr = urp.checkDocumentsForErrors(converted)

	if httpErr != nil {
		return converted, httpErr
	}

	httpErr = urp.callUpdatingDocumentsHandlers(converted, nil)

	if httpErr != nil {
		return converted, httpErr
	}

	httpErr = urp.resource.DatabaseDriver.Update(
		urp.resource.DbModel,
		converted)

	httpErr = urp.callUpdatedDocumentsHandlers(converted, httpErr)

	if httpErr != nil {
		return converted, httpErr
	}

	return converted, nil
}

func (urp *UpdateRequestProcessor) optimizeFields(decodedJsonSlice []map[string]any) {
	optimizeOnUpdateFields := urp.resource.JSONSchemaConfig.OptimizeOnUpdateFields
	optimizeOnUpdateAnyField := funk.ContainsString(optimizeOnUpdateFields, FIELD_ANY)

	for _, jsonObject := range urp.request.DecodedJsonSlice {
		urp.preRequestOptimizeFields(
			jsonObject,
			optimizeOnUpdateFields,
			optimizeOnUpdateAnyField)
	}
}

func (urp *UpdateRequestProcessor) preRequestJSONActions(jsonDocuments []map[string]any) {
	var httpErr HTTPError

	requireOnUpdateFields := urp.resource.JSONSchemaConfig.RequiredOnUpdateFields
	updatableFields := urp.resource.JSONSchemaConfig.UpdatableFields

	requireOnUpdateAnyField := funk.ContainsString(requireOnUpdateFields, FIELD_ANY)
	updatableAnyField := funk.ContainsString(updatableFields, FIELD_ANY)

	if requireOnUpdateAnyField {
		requireOnUpdateFields = urp.resource.DbModelJSONFields
	}

	if updatableAnyField {
		updatableFields = urp.resource.DbModelJSONFields
	}

	for _, jsonObject := range jsonDocuments {
		if httpErr = urp.preRequestRequireOnUpdateChecks(
			jsonObject,
			requireOnUpdateFields); httpErr != nil {
			jsonObject["__http_error__"] = httpErr
			continue
		}

		if httpErr = urp.preRequestUpdatableChecks(
			jsonObject,
			requireOnUpdateFields,
			updatableFields); httpErr != nil {
			jsonObject["__http_error__"] = httpErr
			continue
		}

		if httpErr = urp.preRequestValidateJSON(jsonObject); httpErr != nil {
			jsonObject["__http_error__"] = httpErr
			continue
		}
	}
}

func (urp *UpdateRequestProcessor) checkRanges() HTTPError {
	if urp.request.ContentLength > urp.resource.UpdateMaxInputPayloadSize {
		return NewPayloadTooBigHTTPError(urp.resource.UpdateMaxInputPayloadSize, nil)
	}

	if int64(len(urp.request.Body)) > urp.resource.UpdateMaxInputPayloadSize {
		return NewPayloadTooBigHTTPError(urp.resource.UpdateMaxInputPayloadSize, nil)
	}

	lenDecodedJsonSlice := len(urp.request.DecodedJsonSlice)

	if lenDecodedJsonSlice > int(urp.resource.UpdateMaxInputItems) {
		return NewTooManyInputItemsHTTPError(urp.resource.UpdateMaxInputItems, lenDecodedJsonSlice, nil)
	}

	return nil
}
