package go_cake

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

// TODO require _id and _etag in the payload
// when they are set in the resource
func (drp *DeleteRequestProcessor) ProcessRequest(response *ResponseJSON) ([]GoKateModel, HTTPError) {
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

	httpErr = drp.resource.DatabaseDriver.Delete(
		drp.resource.DbModel,
		converted)

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
