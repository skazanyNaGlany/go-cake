package go_cake

type DatabaseDriver interface {
	TestModel(
		idField string,
		etagField string,
		model GoKateModel,
		dbPath string) error

	Find(
		model GoKateModel,
		where, sort string,
		page, perPage int64) ([]GoKateModel, HTTPError)

	Delete(
		model GoKateModel,
		documents []GoKateModel,
	) HTTPError

	Total(
		model GoKateModel,
		where string) (uint64, HTTPError)

	Insert(
		model GoKateModel,
		documents []GoKateModel,
	) HTTPError

	Update(
		model GoKateModel,
		documents []GoKateModel,
	) HTTPError

	GetWhereFields(where string) ([]string, HTTPError)
	GetSortFields(sort string) ([]string, HTTPError)
	GetProjectionFields(projection string) (map[string]bool, HTTPError)
}

type GoKateModel interface {
	CreateInstance() GoKateModel
	CreateETag() any
	SetID(id string) error
	GetID() any
	SetETag(etag string) error
	GetETag() any
	SetHTTPError(httpError HTTPError)
	GetHTTPError() HTTPError
}

type JSONValidator interface {
	Validate(item map[string]any) error
}

type RequestProcessor interface {
	ProcessRequest(response *ResponseJSON) ([]GoKateModel, HTTPError)
}
