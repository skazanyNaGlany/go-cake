package go_cake

type GoCakeModel interface {
	CreateInstance() GoCakeModel
	ToMap() (map[string]any, error)
	CreateETag() any
	SetID(id string) error
	GetID() any
	SetETag(etag string) error
	GetETag() any
	SetHTTPError(httpError HTTPError)
	GetHTTPError() HTTPError
}
