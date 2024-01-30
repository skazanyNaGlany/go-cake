package go_cake

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	NotFoundHandler http.Handler
	resources       map[string]*Resource
	middlewares     []MiddlewareFunc
}

func NewHandler() *Handler {
	handler := Handler{}
	handler.resources = make(map[string]*Resource)

	return &handler
}

func (rh *Handler) Use(mwf ...MiddlewareFunc) {
	rh.middlewares = append(rh.middlewares, mwf...)
}

func (rh *Handler) processRequest(
	request *Request,
	resource *Resource,
	response *ResponseJSON) {
	if request.IsGet {
		processor := NewGetRequestProcessor(request, resource)

		processor.BaseRequestProcessor.ProcessRequest(response)
	} else if request.IsDelete {
		processor := NewDeleteRequestProcessor(request, resource)

		processor.BaseRequestProcessor.ProcessRequest(response)
	} else if request.IsInsert {
		processor := NewInsertRequestProcessor(request, resource)

		processor.BaseRequestProcessor.ProcessRequest(response)
	} else if request.IsUpdate {
		processor := NewUpdateRequestProcessor(request, resource)

		processor.BaseRequestProcessor.ProcessRequest(response)
	} else if request.IsCORS {
		processor := NewCORSRequestProcessor(request, resource)

		processor.BaseRequestProcessor.ProcessRequest(response)
	} else {
		httpErr := NewMethodNotAllowedHTTPError(nil)

		response.Meta.StatusMessage = httpErr.GetStatusMessage()
		response.Meta.StatusCode = httpErr.GetStatusCode()
	}
}

func (rh *Handler) writeResponse(response *ResponseJSON, httpWriter http.ResponseWriter) {
	jsonText, _ := json.Marshal(response)

	httpWriter.Header().Set("X-GO-KATE-REQUEST-UNIQUE-ID", response.Meta.RequestUniqueID)
	httpWriter.Header().Set("X-GO-KATE-VERSION", response.Meta.Version)
	httpWriter.Header().Set("Content-Type", RESPONSE_CONTENT_TYPE)
	httpWriter.Header().Set("Cache-Control", RESPONSE_CACHE_CONTROL)

	httpWriter.WriteHeader(int(response.Meta.StatusCode))
	httpWriter.Write(jsonText)
}

func (rh *Handler) mainResourceHandler(httpWriter http.ResponseWriter, httpRequest *http.Request) {
	var httpErr HTTPError

	response := NewResponseJSON()

	resource := rh.FindMatchedResource(httpRequest)

	if resource == nil {
		if rh.NotFoundHandler != nil {
			rh.NotFoundHandler.ServeHTTP(httpWriter, httpRequest)

			return
		}

		httpErr = NewURLNotFoundHTTPError(nil)

		response.Meta.StatusMessage = httpErr.GetStatusMessage()
		response.Meta.StatusCode = httpErr.GetStatusCode()

		rh.writeResponse(response, httpWriter)
		return
	}

	request := Request{
		ResourcePattern: resource.CompiledPattern,
		Resource:        resource.ResourceName,
		Method:          httpRequest.Method,
		Request:         httpRequest,
		ResponseWriter:  httpWriter,
	}

	if httpErr = request.Parse(httpRequest); httpErr != nil {
		response.Meta.StatusMessage = httpErr.GetStatusMessage()
		response.Meta.StatusCode = httpErr.GetStatusCode()

		rh.writeResponse(response, request.ResponseWriter)
		return
	}

	rh.processRequest(&request, resource, response)

	if request.ResponseWriter != nil {
		rh.writeResponse(response, request.ResponseWriter)
	}
}

func (rh *Handler) targetHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rh.mainResourceHandler(w, r)
	})
}

func (rh *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handlers := rh.middlewares
	handlers = append(handlers, rh.targetHandler)

	queue := MiddlewareQueue{
		ResponseWriter: w,
		Request:        r,
		Queue:          handlers,
	}

	queue.Execute()
}

func (rh *Handler) FindMatchedResource(r *http.Request) *Resource {
	for _, iresource := range rh.resources {
		if iresource.MatchPattern(r.URL.Path) {
			return iresource
		}
	}

	return nil
}

func (rh *Handler) AddResource(resource *Resource) error {
	if _, exists := rh.resources[resource.Pattern]; exists {
		return &RestHandlerPatternExistsError{}
	}

	rh.resources[resource.Pattern] = resource

	return nil
}
