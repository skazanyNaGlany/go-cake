package go_cake

type GetRequestProcessor struct {
	BaseRequestProcessor
}

func NewGetRequestProcessor(request *Request, resource *Resource) *GetRequestProcessor {
	var getRequestProcessor GetRequestProcessor

	getRequestProcessor.request = request
	getRequestProcessor.resource = resource
	getRequestProcessor.subRequestProcessor = &getRequestProcessor

	return &getRequestProcessor
}

func (grp *GetRequestProcessor) ProcessRequest(response *ResponseJSON) ([]GoKateModel, HTTPError) {
	var httpErr HTTPError
	var documents []GoKateModel

	if !grp.resource.GetAllowed {
		return nil, NewMethodNotAllowedHTTPError(nil)
	}

	grp.initPagination(response)

	if httpErr = grp.checkRanges(); httpErr != nil {
		return nil, httpErr
	}

	if httpErr = grp.preRequestModelActions(); httpErr != nil {
		return nil, httpErr
	}

	if grp.request.PerPage == 0 {
		return nil, nil
	}

	documents, httpErr = grp.resource.DatabaseDriver.Find(
		grp.resource.DbModel,
		grp.request.Where,
		grp.request.Sort,
		grp.request.Page,
		grp.request.PerPage)

	httpErr = grp.callFetchedDocumentsHandlers(documents, httpErr)

	if httpErr != nil {
		return nil, httpErr
	}

	return documents, nil
}

func (grp *GetRequestProcessor) preRequestModelActions() HTTPError {
	// no actions for get
	return nil
}

func (grp *GetRequestProcessor) initPagination(response *ResponseJSON) {
	if grp.request.PerPage == 0 {
		grp.request.PerPage = grp.resource.GetMaxOutputItems
	}

	response.Meta.Page = grp.request.Page
	response.Meta.PerPage = grp.request.PerPage
}

func (grp *GetRequestProcessor) checkRanges() HTTPError {
	if grp.request.PerPage > grp.resource.GetMaxOutputItems {
		return NewPerPageTooLargeHTTPError(grp.resource.GetMaxOutputItems, nil)
	}

	return nil
}
