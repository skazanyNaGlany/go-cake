package go_cake

import (
	"github.com/thoas/go-funk"
)

type InsertRequestProcessor struct {
	BaseRequestProcessor
}

func NewInsertRequestProcessor(request *Request, resource *Resource) *InsertRequestProcessor {
	var insertRequestProcessor InsertRequestProcessor

	insertRequestProcessor.request = request
	insertRequestProcessor.resource = resource
	insertRequestProcessor.subRequestProcessor = &insertRequestProcessor

	return &insertRequestProcessor
}

func (irp *InsertRequestProcessor) ProcessRequest(response *ResponseJSON) ([]GoKateModel, HTTPError) {
	var httpErr HTTPError

	if !irp.resource.InsertAllowed {
		return nil, NewMethodNotAllowedHTTPError(nil)
	}

	if httpErr = irp.checkRanges(); httpErr != nil {
		return nil, httpErr
	}

	if irp.request.HasWhere() || irp.request.HasSort() || irp.request.HasPage() {
		return nil, NewModifiersNotAllowedHTTPError(nil)
	}

	irp.optimizeFields(irp.request.DecodedJsonSlice)
	irp.preRequestJSONActions(irp.request.DecodedJsonSlice)

	converted, err := irp.decodedJsonSliceToDBModels(irp.request.DecodedJsonSlice)

	if err != nil {
		return converted, err
	}

	httpErr = irp.checkDocumentsForErrors(converted)

	if httpErr != nil {
		return converted, httpErr
	}

	httpErr = irp.callInsertingDocumentsHandlers(converted, nil)

	if httpErr != nil {
		return converted, httpErr
	}

	httpErr = irp.resource.DatabaseDriver.Insert(
		irp.resource.DbModel,
		converted)

	httpErr = irp.callInsertedDocumentsHandlers(converted, httpErr)

	if httpErr != nil {
		return converted, httpErr
	}

	return converted, nil
}

func (irp *InsertRequestProcessor) optimizeFields(decodedJsonSlice []map[string]any) {
	optimizeOnInsertFields := irp.resource.JSONSchemaConfig.OptimizeOnInsertFields
	optimizeOnInsertAnyField := funk.ContainsString(optimizeOnInsertFields, FIELD_ANY)

	for _, jsonObject := range irp.request.DecodedJsonSlice {
		irp.preRequestOptimizeFields(
			jsonObject,
			optimizeOnInsertFields,
			optimizeOnInsertAnyField)
	}
}

func (irp *InsertRequestProcessor) preRequestJSONActions(jsonDocuments []map[string]any) {
	var httpErr HTTPError

	requireOnInsertFields := irp.resource.JSONSchemaConfig.RequiredOnInsertFields
	insertableFields := irp.resource.JSONSchemaConfig.InsertableFields

	requireOnInsertAnyField := funk.ContainsString(requireOnInsertFields, FIELD_ANY)
	insertableAnyField := funk.ContainsString(insertableFields, FIELD_ANY)

	if requireOnInsertAnyField {
		requireOnInsertFields = irp.resource.DbModelJSONFieldsNoReserved
	}

	if insertableAnyField {
		insertableFields = irp.resource.DbModelJSONFieldsNoReserved
	}

	for _, jsonObject := range irp.request.DecodedJsonSlice {
		if httpErr = irp.preRequestRequireOnInsertChecks(
			jsonObject,
			requireOnInsertFields); httpErr != nil {
			jsonObject["__http_error__"] = httpErr
			continue
		}

		if httpErr = irp.preRequestInsertableChecks(
			jsonObject,
			requireOnInsertFields,
			insertableFields); httpErr != nil {
			jsonObject["__http_error__"] = httpErr
			continue
		}

		if httpErr = irp.preRequestValidateJSON(jsonObject); httpErr != nil {
			jsonObject["__http_error__"] = httpErr
			continue
		}
	}
}

func (irp *InsertRequestProcessor) checkRanges() HTTPError {
	if irp.request.ContentLength > irp.resource.InsertMaxInputPayloadSize {
		return NewPayloadTooBigHTTPError(irp.resource.InsertMaxInputPayloadSize, nil)
	}

	if int64(len(irp.request.Body)) > irp.resource.InsertMaxInputPayloadSize {
		return NewPayloadTooBigHTTPError(irp.resource.InsertMaxInputPayloadSize, nil)
	}

	lenDecodedJsonSlice := len(irp.request.DecodedJsonSlice)

	if lenDecodedJsonSlice > int(irp.resource.InsertMaxInputItems) {
		return NewTooManyInputItemsHTTPError(irp.resource.InsertMaxInputItems, lenDecodedJsonSlice, nil)
	}

	return nil
}
