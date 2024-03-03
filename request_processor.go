package go_cake

type RequestProcessor interface {
	ProcessRequest(response *ResponseJSON) ([]GoCakeModel, HTTPError)
}
