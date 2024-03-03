package go_cake

type DocumentsCallback func(
	resource *Resource,
	request *Request,
	documents []GoCakeModel,
	currentHttpErr HTTPError) HTTPError
