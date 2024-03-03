package go_cake

type PrePostRequestCallback func(
	resource *Resource,
	request *Request,
	response *ResponseJSON) HTTPError
