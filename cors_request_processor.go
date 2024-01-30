package go_cake

type CORSRequestProcessor struct {
	BaseRequestProcessor
}

func NewCORSRequestProcessor(request *Request, resource *Resource) *CORSRequestProcessor {
	var corsRequestProcessor CORSRequestProcessor

	corsRequestProcessor.request = request
	corsRequestProcessor.resource = resource
	corsRequestProcessor.subRequestProcessor = &corsRequestProcessor

	return &corsRequestProcessor
}

func (crp *CORSRequestProcessor) ProcessRequest(response *ResponseJSON) ([]GoKateModel, HTTPError) {
	return nil, nil
}
