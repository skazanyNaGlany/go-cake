package go_cake

import "context"

type DatabaseDriver interface {
	TestModel(
		idField string,
		etagField string,
		model GoKateModel,
		dbPath string) error

	Find(
		model GoKateModel,
		where, sort string,
		page, perPage int64,
		ctx context.Context,
		userData any) ([]GoKateModel, HTTPError)

	Delete(
		model GoKateModel,
		documents []GoKateModel,
		ctx context.Context,
		userData any) HTTPError

	Total(
		model GoKateModel,
		where string,
		ctx context.Context,
		userData any) (uint64, HTTPError)

	Insert(
		model GoKateModel,
		documents []GoKateModel,
		ctx context.Context,
		userData any) HTTPError

	Update(
		model GoKateModel,
		documents []GoKateModel,
		ctx context.Context,
		userData any) HTTPError

	GetWhereFields(model GoKateModel, where string) ([]string, HTTPError)
	GetSortFields(model GoKateModel, sort string) ([]string, HTTPError)
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
