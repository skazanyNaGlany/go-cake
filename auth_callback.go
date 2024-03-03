package go_cake

// app callback
type AuthCallback func(
	resource *Resource,
	request *Request,
	response *ResponseJSON) bool
