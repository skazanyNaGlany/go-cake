package go_cake

import "context"

type DatabaseDriver interface {
	TestModel(
		idField string,
		etagField string,
		model GoCakeModel,
		dbPath string) error

	Find(
		model GoCakeModel,
		where, sort string,
		page, perPage int64,
		ctx context.Context,
		userData any) ([]GoCakeModel, HTTPError)

	Delete(
		model GoCakeModel,
		documents []GoCakeModel,
		ctx context.Context,
		userData any) HTTPError

	Total(
		model GoCakeModel,
		where string,
		ctx context.Context,
		userData any) (uint64, HTTPError)

	Insert(
		model GoCakeModel,
		documents []GoCakeModel,
		ctx context.Context,
		userData any) HTTPError

	Update(
		model GoCakeModel,
		documents []GoCakeModel,
		ctx context.Context,
		userData any) HTTPError

	GetWhereFields(model GoCakeModel, where string) ([]string, HTTPError)
	GetSortFields(model GoCakeModel, sort string) ([]string, HTTPError)
}

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

type JSONValidator interface {
	Validate(item map[string]any) error
}

type RequestProcessor interface {
	ProcessRequest(response *ResponseJSON) ([]GoCakeModel, HTTPError)
}
