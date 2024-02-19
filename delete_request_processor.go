package go_cake

import "github.com/thoas/go-funk"

type DeleteRequestProcessor struct {
	BaseRequestProcessor
}

func NewDeleteRequestProcessor(request *Request, resource *Resource) *DeleteRequestProcessor {
	var deleteRequestProcessor DeleteRequestProcessor

	deleteRequestProcessor.request = request
	deleteRequestProcessor.resource = resource
	deleteRequestProcessor.subRequestProcessor = &deleteRequestProcessor

	return &deleteRequestProcessor
}

func (drp *DeleteRequestProcessor) ProcessRequest(response *ResponseJSON) ([]GoCakeModel, HTTPError) {
	var httpErr HTTPError

	if !drp.resource.DeleteAllowed {
		return nil, NewMethodNotAllowedHTTPError(nil)
	}

	if httpErr = drp.checkRanges(); httpErr != nil {
		return nil, httpErr
	}

	if drp.request.HasWhere() || drp.request.HasSort() || drp.request.HasPage() {
		return nil, NewModifiersNotAllowedHTTPError(nil)
	}

	drp.optimizeFields(drp.request.DecodedJsonSlice)
	drp.preRequestJSONActions(drp.request.DecodedJsonSlice)

	converted, err := drp.decodedJsonSliceToDBModels(drp.request.DecodedJsonSlice)

	if err != nil {
		return converted, err
	}

	httpErr = drp.checkDocumentsForErrors(converted)

	if httpErr != nil {
		return converted, httpErr
	}

	httpErr = drp.callDeletingDocumentsHandlers(converted, nil)

	if httpErr != nil {
		return converted, httpErr
	}

	ctx, cancel := drp.resource.ResourceCallback.CreateContext(
		drp.resource,
		drp.request,
		response,
		ctxDbDriverDelete)
	defer cancel()

	httpErr = drp.resource.DatabaseDriver.Delete(
		drp.resource.DbModel,
		converted,
		ctx,
		nil)

	httpErr = drp.callDeletedDocumentsHandlers(converted, httpErr)

	if httpErr != nil {
		return converted, httpErr
	}

	return converted, nil
}

func (drp *DeleteRequestProcessor) checkRanges() HTTPError {
	if drp.request.ContentLength > drp.resource.DeleteMaxInputPayloadSize {
		return NewPayloadTooBigHTTPError(drp.resource.DeleteMaxInputPayloadSize, nil)
	}

	if int64(len(drp.request.Body)) > drp.resource.DeleteMaxInputPayloadSize {
		return NewPayloadTooBigHTTPError(drp.resource.DeleteMaxInputPayloadSize, nil)
	}

	lenDecodedJsonSlice := len(drp.request.DecodedJsonSlice)

	if lenDecodedJsonSlice > int(drp.resource.DeleteMaxInputItems) {
		return NewTooManyInputItemsHTTPError(drp.resource.DeleteMaxInputItems, lenDecodedJsonSlice, nil)
	}

	return nil
}

func (drp *DeleteRequestProcessor) optimizeFields(decodedJsonSlice []map[string]any) {
	optimizeOnDeleteFields := drp.resource.JSONSchemaConfig.OptimizeOnDeleteFields
	optimizeOnDeleteAnyField := funk.ContainsString(optimizeOnDeleteFields, FIELD_ANY)

	for _, jsonObject := range drp.request.DecodedJsonSlice {
		drp.preRequestOptimizeFields(
			jsonObject,
			optimizeOnDeleteFields,
			optimizeOnDeleteAnyField)
	}
}

func (drp *DeleteRequestProcessor) preRequestJSONActions(jsonDocuments []map[string]any) {
	var httpErr HTTPError

	requireOnDeleteFields := drp.resource.JSONSchemaConfig.RequiredOnDeleteFields
	requireOnDeleteAnyField := funk.ContainsString(requireOnDeleteFields, FIELD_ANY)

	if requireOnDeleteAnyField {
		requireOnDeleteFields = drp.resource.DbModelJSONFields
	}

	for _, jsonObject := range jsonDocuments {
		if httpErr = drp.preRequestRequireOnUpdateChecks(
			jsonObject,
			requireOnDeleteFields); httpErr != nil {
			jsonObject["__http_error__"] = httpErr
			continue
		}

		if httpErr = drp.preRequestValidateJSON(jsonObject); httpErr != nil {
			jsonObject["__http_error__"] = httpErr
			continue
		}
	}
}

func (drp *DeleteRequestProcessor) preRequestValidateJSON(
	jsonObjectMap map[string]any) HTTPError {
	if drp.resource.JSONSchemaConfig == nil ||
		drp.resource.JSONSchemaConfig.DeleteValidator == nil {
		return nil
	}

	if err := drp.resource.JSONSchemaConfig.DeleteValidator.Validate(jsonObjectMap); err != nil {
		return NewClientObjectMalformedHTTPError(err)
	}

	return nil
}
