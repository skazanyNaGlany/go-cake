package go_cake

import "context"

type CreateContextCallback func(
	resource *Resource,
	request *Request,
	response *ResponseJSON,
	contextType ContextType) (context.Context, context.CancelFunc)
